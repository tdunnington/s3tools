/* Copyright 2015 Timothy Eric Dunnington
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 // lib package contains core s3 functions for the cmdline executables -
 // s3 location format parsers, basic file functions for removing, uploading,
 // downloading and listing, etc.
package lib

import (
  "github.com/tdunnington/s3tools/logger"
  "fmt"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/s3"
  "os"
  "regexp"
)

// s3pathre is the regular expression for parsing s3 paths
const s3pathre = "^s3:([^:]+):(.+)$"

// S3path is a struct representing the parts of the s3 path
type S3path struct {
  Bucket string
  Path string
}

// Download a file from an S3 location identified by "source", to a local file
// identified by "destination". Returns error if the source path is not in
// the proper format, if the source path is not readable with the current
// credentials, or if the destination is not writeable.
func CopyFromS3(source, destination string) error {
  src, err := ParseS3Path(source)
	if err != nil {
		return fmt.Errorf("Source path invalid, error was: %s\n", err)
	}

	writer, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("Failed to open file '%s', error was: %s\n", destination, err)
	}

	downloader := s3manager.NewDownloader(nil)
	if downloader == nil {
		return fmt.Errorf("Internal Error: Failure creating NewDownloader @copyToS3()\n")
	}

	defer writer.Close()
	bytesWritten, err := downloader.Download(writer, &s3.GetObjectInput{
		Bucket:   aws.String(src.Bucket),
		Key:      aws.String(src.Path),
	})

	if err != nil {
		return fmt.Errorf("Failed to download source file '%s' to destination '%s'\nError from S3 was: %s\n", source, destination, err)
	}

	logger.Debug(fmt.Sprintf("Downloaded '%s', %d bytes retrieved\n", source, bytesWritten))

	return nil
}

// Upload a file to an S3 location identified by "destination",
// from a local file identified by "source". Returns an error if the source
// file is not readable, if the destination s3 path is not writable, or if
// the destination path is not in the right format.
func CopyToS3(source string, destination string, isReducedRedundancy bool) error {
  dest, err := ParseS3Path(destination)
	if err != nil {
		return fmt.Errorf("Destination path invalid, error was: %s\n", err)
	}

	reader, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Failed to open file '%s', error was: %s\n", source, err)
	}

	uploader := s3manager.NewUploader(nil)
	if uploader == nil {
		return fmt.Errorf("Internal Error: Failure creating NewUploader @copyToS3()\n")
	}

	uploadInput := s3manager.UploadInput{
		Body:     reader,
		Bucket:   aws.String(dest.Bucket),
		Key:      aws.String(dest.Path),
	}

	if isReducedRedundancy {
		uploadInput.StorageClass = aws.String("REDUCED_REDUNDANCY")
	}

	defer reader.Close()
	result, err := uploader.Upload(&uploadInput)

	if err != nil {
		return fmt.Errorf("Failed to upload source file '%s' to destination '%s'\nError from S3 was: %s\n", source, destination, err)
	}

	logger.Debug(fmt.Sprintf("Post-upload file destination URL:='%s'\n",result.Location))

	return nil
}

// RemoveS3Path removes an object in an S3 location specified in "path".
// Returns an error if the "path" is not a valid s3 path, or if the
// service failed to remove the object
func RemoveS3Path(path string) error {
  target, err := ParseS3Path(path)
  if err != nil {
    return err
  }

  svc := s3.New(nil)
  params := &s3.DeleteObjectInput{
      Bucket:       aws.String(target.Bucket),
      Key:          aws.String(target.Path),
  }
  result, err := svc.DeleteObject(params)

  if err != nil {
      return err
  }

  logger.Debug(fmt.Sprintf("Removed item ''%s'\n",result))
  return nil
}

// IsS3Path returns true if the string s is a valid s3 path, false otherwise
func IsS3Path(s string) bool {
  isMatch, _ := regexp.MatchString(s3pathre,s)
  return isMatch
}

// ParseS3Path returns a struct of type S3path, from the string "path". Returns
// an error if the "path" is not a valid s3 path of the form
// "s3:bucket:/path/to/file"
func ParseS3Path(path string) (*S3path, error) {
  s := new(S3path)
  re := regexp.MustCompile(s3pathre)
	parts := re.FindStringSubmatch(path)
	s.Bucket = parts[1]
	s.Path = parts[2]

  if len(s.Bucket) == 0 || len(s.Path) == 0 {
    return nil, fmt.Errorf("The path '%s' is invalid, must be in the form s3:bucket:/path/to/file\n", path)
  }

  return s, nil
}
