/** 
 * Copyright (c) 2017, Patrick Uiterwijk <patrick@puiterwijk.org>
 * All rights reserved.
 *
 * This file is part of GoBackLight.
 *
 * GoBackLight is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GoBackLight is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GoBackLight.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "os"
)

func findClassName() string {
    return "intel_backlight"
}

func getFileName(name string) string {
    return "/sys/class/backlight/" + findClassName() + "/" + name
}

func getValue(name string) int {
    rawValue, err := ioutil.ReadFile(getFileName(name))
    if err != nil {
        panic(err)
    }
    strValue := string(rawValue)
    strValue = strings.TrimSpace(strValue)
    value, err := strconv.Atoi(strValue)
    if err != nil {
        panic(err)
    }
    return value
}

func setValue(name string, value int) {
    strValue := strconv.Itoa(value)
    rawValue := []byte(strValue)
    err := ioutil.WriteFile(getFileName(name), rawValue, 0644)
    if err != nil {
        panic(err)
    }
}

type Operation uint8
const (
    OP_GET Operation = iota + 1
    OP_SET_ABS
    OP_SET_PERC
    OP_ABORT
)

type Arguments struct {
    operation Operation
    amount int
    verbose bool
    relative bool
}

func getArguments() Arguments {
    outp := Arguments{OP_GET, 0, false, false}
    outp.amount = 1
    if len(os.Args) == 1 {
        return outp
    } else {
        nextArg := os.Args[1]
        if nextArg == "-v" {
            outp.verbose = true
            if len(os.Args) == 2 {
                return outp
            }
            outp.operation = OP_SET_ABS
            nextArg = os.Args[2]
        }
        if strings.HasPrefix(nextArg, "+") {
            outp.relative = true
            nextArg = nextArg[1:]
        } else if strings.HasPrefix(nextArg, "-") {
            outp.relative = true
            nextArg = nextArg[1:]
            outp.amount = -1
        }
        if strings.HasSuffix(nextArg, "%") {
            outp.operation = OP_SET_PERC
            nextArg = strings.TrimSuffix(nextArg, "%")
        }
        amount, err := strconv.Atoi(nextArg)
        if err != nil {
            fmt.Println("Unable to parse value")
            outp.operation = OP_ABORT
        } else {
            outp.amount = outp.amount * amount
        }
        return outp
    }
}

func main() {
    maxBrightness := getValue("max_brightness")
    currentBrightness := getValue("brightness")

    arguments := getArguments()

    if arguments.operation == OP_GET {
        fmt.Println(strconv.Itoa(currentBrightness) + " / " + strconv.Itoa(maxBrightness))
    } else if arguments.operation == OP_SET_ABS {
        newBrightness := arguments.amount
        if arguments.relative {
            newBrightness += currentBrightness
        }
        if newBrightness > maxBrightness {
            if arguments.verbose {
                fmt.Println("Capping to max brightness")
            }
            newBrightness = maxBrightness
        }
        if newBrightness < 0 {
            if arguments.verbose {
                fmt.Println("Capping to min brightness")
            }
            newBrightness = 0
        }
        if arguments.verbose {
            fmt.Println(strconv.Itoa(currentBrightness) + " => " + strconv.Itoa(newBrightness))
        }
        setValue("brightness", newBrightness)
        if arguments.verbose {
            fmt.Println("Written")
        }
    } else if arguments.operation == OP_SET_PERC {
        if arguments.amount < -100 || arguments.amount > 100 {
            fmt.Println("Percentage out of range")
            return
        }
        newBrightness := int((float32(maxBrightness) / 100.0) * float32(arguments.amount))
        if arguments.relative {
            newBrightness = currentBrightness + newBrightness
        }
        if newBrightness > maxBrightness {
            if arguments.verbose {
                fmt.Println("Capping to max brightness")
            }
            newBrightness = maxBrightness
        }
        if newBrightness < 0 {
            if arguments.verbose {
                fmt.Println("Capping to min brightness")
            }
            newBrightness = 0
        }
        if arguments.verbose {
            fmt.Println(strconv.Itoa(currentBrightness) + " => " + strconv.Itoa(newBrightness))
        }
        setValue("brightness", newBrightness)
        if arguments.verbose {
            fmt.Println("Written")
        }
    } else {
        if arguments.verbose {
            fmt.Println("Doing nothing")
        }
    }
}
