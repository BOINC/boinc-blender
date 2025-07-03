// This file is part of BOINC.
// https://boinc.berkeley.edu
// Copyright (C) 2025 University of California
//
// BOINC is free software; you can redistribute it and/or modify it
// under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation,
// either version 3 of the License, or (at your option) any later version.
//
// BOINC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with BOINC.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
)

func parseProgress(line string) (current, total int, found bool) {
	re := regexp.MustCompile(` (\d+)\s?/\s?(\d+)$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		if curr, err := strconv.Atoi(matches[1]); err == nil {
			if tot, err := strconv.Atoi(matches[2]); err == nil {
				return curr, tot, true
			}
		}
	}
	return 0, 0, false
}

func writeProgress(progress float64) {
	progressFile, err := os.OpenFile("fraction_done", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening progress file: %v\n", err)
		return
	}
	defer progressFile.Close()

	if _, err := progressFile.WriteString(fmt.Sprintf("%.3f", progress)); err != nil {
		fmt.Printf("Error writing to progress file: %v\n", err)
	}
}

func main() {
	const VERSION = "1.0.0"
	var (
		blenderPath  = flag.String("blender", "/bin/blender", "Path to Blender executable")
		workDir      = flag.String("workdir", "/app", "Working directory for Blender")
		output       = flag.String("output", "", "Prefix for rendered images")
		inputFile    = flag.String("input", "", "Input Blender file (.blend)")
		frame        = flag.Int("frame", 0, "Frame to render")
		engine       = flag.String("engine", "CYCLES", "Render engine (CYCLES, EEVEE)")
		cyclesDevice = flag.String("cyclesDevice", "CPU", "Cycles device to use (CPU, CUDA, OPTIX, HIP, ONEAPI, METAL)")
		version      = flag.Bool("version", false, "Show version information")
		help         = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help {
		fmt.Printf("BOINC Blender Application v%s\n", VERSION)
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	if *version {
		fmt.Printf("BOINC Blender Application v%s\n", VERSION)
		return
	}

	if *inputFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("BOINC Blender Application v%s\n", VERSION)
	fmt.Printf("Blender Path: %s\n", *blenderPath)
	fmt.Printf("Working directory: %s\n", *workDir)
	fmt.Printf("Output directory: %s\n", *output)
	fmt.Printf("Input file: %s\n", *inputFile)
	fmt.Printf("Frame to render: %d\n", *frame)
	fmt.Printf("Render engine: %s\n", *engine)
	fmt.Printf("Cycles device: %s\n", *cyclesDevice)

	var renderEngine = ""
	var renderEngineArg = ""
	switch *engine {
	case "CYCLES":
		renderEngine = "CYCLES"
		if *cyclesDevice != "" {
			devices := []string{"CPU", "CUDA", "OPTIX", "HIP", "ONEAPI", "METAL"}
			if !slices.Contains(devices, *cyclesDevice) {
				fmt.Println("Unsupported Cycles device. Please use CPU, CUDA, OPTIX, HIP, ONEAPI, or METAL.")
				os.Exit(1)
			}
			renderEngineArg = fmt.Sprintf("-- --cycles-device %s", *cyclesDevice)
		}
	case "EEVEE":
		renderEngine = "BLENDER_EEVEE_NEXT"
	default:
		fmt.Println("Unsupported render engine. Please use CYCLES or EEVEE.")
		os.Exit(1)
	}

	cmd := fmt.Sprintf("%s -b %s/%s -o %s/%s -F PNG -f %d -t 1 -E %s --factory-startup %s",
		*blenderPath,
		*workDir,
		*inputFile,
		*workDir,
		*output,
		*frame,
		renderEngine,
		renderEngineArg,
	)
	fmt.Println("Executing command:", cmd)
	blenderCmd := exec.Command("sh", "-c", cmd)

	stdout, err := blenderCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	stderr, err := blenderCmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		os.Exit(1)
	}

	if err := blenderCmd.Start(); err != nil {
		fmt.Printf("Error starting Blender command: %v\n", err)
		os.Exit(1)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			if current, total, found := parseProgress(line); found {
				writeProgress(float64(current) / float64(total))
			}
			fmt.Printf("%s\n", line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stdout: %v\n", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%s\n", line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stderr: %v\n", err)
		}
	}()

	if err := blenderCmd.Wait(); err != nil {
		fmt.Printf("Error running Blender command: %v\n", err)
		os.Exit(1)
	}
	writeProgress(1.0)
	fmt.Println("Blender command executed successfully.")
}
