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

// logger is a package of logging functions for the s3tools project
package logger

import (
  "fmt"
)

var IsDebugMode bool = false
var IsQuietMode bool = false

// Log writes the string "s" to stdout if IsQuietMode is false, or if
// IsDebugMode is true. Otherwise, it does nothing.
func Log(s string) {
  // this ensures that debug mode overrides
  // quiet mode
  if !IsQuietMode || IsDebugMode {
    fmt.Printf("%s\n", s)
  }
}

// Debug writes the string "s" to stdout if IsDebugMode is true, otherwise
// it does nothing.
func Debug(s string) {
  if IsDebugMode {
    fmt.Printf("DEBUG: %s\n", s)
  }
}
