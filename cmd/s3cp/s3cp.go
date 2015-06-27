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

/* s3cp
 *
 * Uploads or downloads a file from an S3 bucket, using scp conventions.
 *
 * USAGE:
 *
 *     s3cp [--help] [--quiet] [--debug] [--region regionname] [--rr] source destination
 *
 * WHERE:
 *
 *     --help          : prints help
 *
 *     --quiet         : suppress output
 *
 *     --debug         : debug mode, prints lots of debug info
 *
 *     --region        : the AWS region of the target bucket; defaults to us-east-1
 *
 *     --rr            : sets the upload to reduced redundancy; has no effect on download
 *
 *     "source" and "desination" can be either local or remote objects, and both
 *     are required.
 *
 *     For a remote object: s3:bucket:/folder.../file.name
 *
 *     For a local object : /folder.../file.name
 *
 * EXAMPLES:
 *
 *     s3cp s3:/mybucket/myfolder/backup.tar.gz /tmp
 *       - downloads backup.tar.gz from S3 and places it in /tmp folder
 *
 *     s3cp s3:/mybucket/myfolder/backup.tar.gz /tmp/foobar.tar.gz
 *       - downloads backup.tar.gz from S3 and places it in the file /tmp/foobar.tar.gz
 *
 *     s3cp /tmp/backup.tar.gz s3:/mybucket/myfolder
 *       - uploads backup.tar.gz to S3 in the bucket mybucket and folder myfolder
 *
 */

package main

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"flag"
	"github.com/tdunnington/s3tools/logger"
	"github.com/tdunnington/s3tools/lib"
)

var region string = "us-east-1"
var isReducedRedundancy bool = false

// Prints the help information to stdout
func printHelp() {
	// print help using the flag package default
	flag.Usage()
	// add the two trailing arguments for source and dest
	fmt.Fprintf(os.Stderr, "  source:  The source of the copy, either a local file path or an s3 path like s3:bucket:/path\n")
	fmt.Fprintf(os.Stderr, "  destination:  The destination of the copy, in the same format as source (above)\n")
	fmt.Fprintf(os.Stderr, "\nBoth source and destination are required, and one must be an s3 path, another must be a local path\n\n")
}

// Parses the cmdline options and places settings in the proper variables.
// Returns error if the cmdline is invalid or if the user requested help
func parseCmdline() (string, string, error) {
	var isHelpRequested bool = false

	flag.BoolVar(&isHelpRequested, "help", false, "(optional) Prints this help message")
	flag.BoolVar(&logger.IsDebugMode, "debug", false, "(optional) Used for debugging; outputs lots of debug info")
	flag.BoolVar(&logger.IsQuietMode, "quiet", false, "(optional) Suppresses output")
	flag.StringVar(&region, "region", "us-east-1", "(optional) The AWS region holding the target bucket; defaults to 'us-east-1'")
	flag.BoolVar(&isReducedRedundancy, "rr", false, "(optional) Sets the upload to reduced redundancy")
	flag.Parse()

	logger.Debug(fmt.Sprintf("got args '%s'\n", os.Args))

	// tried to use flag.ErrHelp here, but it didn't work as expected...seems that ErrHelp is
	// always being set to the "help requested" value. I wonder if the flag package assumes
	// all flags are required or something
	if flag.NArg() != 2 || isHelpRequested {
		logger.Debug(fmt.Sprintf("NArg = %d, ErrHelp = %s\n", flag.NArg(), flag.ErrHelp))
		return "", "", fmt.Errorf("Invalid cmdline arguments or help requested\n")
	}

	return flag.Arg(0), flag.Arg(1), nil
}

func main() {
	source, destination, err := parseCmdline()
	if err != nil {
		printHelp()
		os.Exit(1)
	}

	aws.DefaultConfig.Region = region

	if lib.IsS3Path(source) {
		logger.Debug("Calling copyFromS3")
		err = lib.CopyFromS3(source,destination)
	} else if lib.IsS3Path(destination) {
		logger.Debug("Calling copyToS3")
		err = lib.CopyToS3(source,destination,isReducedRedundancy)
	} else {
		printHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("%s\n",err.Error())
		os.Exit(1)
	}

	logger.Log(fmt.Sprintf("%s -> %s : transfer complete\n",source,destination))
	os.Exit(0)
}
