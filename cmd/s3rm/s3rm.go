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
 
// IAM Policy Requirement : s3:DeleteObject

package main

import (
  "fmt"
  "os"
  "flag"
  "github.com/tdunnington/s3tools/lib"
  "github.com/tdunnington/s3tools/logger"
  "github.com/aws/aws-sdk-go/aws"
)

var region string = "us-east-1"

// Prints the help information to stdout
func printHelp() {
	// print help using the flag package default
	flag.Usage()
	// add the two trailing arguments for source and dest
	fmt.Fprintf(os.Stderr, "  path:  The S3 object do delete, like s3:bucket:/path\n")
}

// Parses the cmdline options and places settings in the proper variables.
// Returns error if the cmdline is invalid or if the user requested help
func parseCmdline() (string, error) {
	var isHelpRequested bool = false

	flag.BoolVar(&isHelpRequested, "help", false, "(optional) Prints this help message")
	flag.BoolVar(&logger.IsDebugMode, "debug", false, "(optional) Used for debugging; outputs lots of debug info")
	flag.BoolVar(&logger.IsQuietMode, "quiet", false, "(optional) Suppresses output")
	flag.StringVar(&region, "region", "us-east-1", "(optional) The AWS region holding the target bucket; defaults to 'us-east-1'")
	flag.Parse()

	logger.Debug(fmt.Sprintf("got args '%s'\n", os.Args))

	// tried to use flag.ErrHelp here, but it didn't work as expected...seems that ErrHelp is
	// always being set to the "help requested" value. I wonder if the flag package assumes
	// all flags are required or something
	if flag.NArg() != 1 || isHelpRequested {
		logger.Debug(fmt.Sprintf("NArg = %d, ErrHelp = %s\n", flag.NArg(), flag.ErrHelp))
		return "", fmt.Errorf("Invalid cmdline arguments or help requested\n")
	}

	return flag.Arg(0), nil
}

func main() {
  path, err := parseCmdline()
	if err != nil {
		printHelp()
		os.Exit(1)
	}

	aws.DefaultConfig.Region = region

	if lib.IsS3Path(path) {
		logger.Debug("Calling RemoveS3Path")
		err = lib.RemoveS3Path(path)
	} else {
		printHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("%s\n",err.Error())
		os.Exit(1)
	}

	logger.Log(fmt.Sprintf("%s removed\n",path))
	os.Exit(0)
}
