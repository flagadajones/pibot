package main

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	monitors := pixelgl.Monitors()

	for i, m := range monitors {

		fmt.Printf("monitor %v:\n", i)

		name := m.Name()
		fmt.Printf("-name: %v\n", name)

		bitDepthRed, bitDepthGreen, bitDepthBlue := m.BitDepth()
		fmt.Printf("-bitDepth: %v-bit red, %v-bit green, %v-bit blue\n",
			bitDepthRed, bitDepthGreen, bitDepthBlue)

		physicalSizeWidth, physicalSizeHeight := m.PhysicalSize()
		fmt.Printf("-physicalSize: %v mm, %v mm\n",
			physicalSizeWidth, physicalSizeHeight)

		positionX, positionY := m.Position()
		fmt.Printf("-position: %v, %v upper-left corner\n",
			positionX, positionY)

		refreshRate := m.RefreshRate()
		fmt.Printf("-refreshRate: %v Hz\n", refreshRate)

		sizeWidth, sizeHeight := m.Size()
		fmt.Printf("-size: %v px, %v px\n",
			sizeWidth, sizeHeight)

		videoModes := m.VideoModes()

		for j, vm := range videoModes {

			fmt.Printf("-video mode %v: -width: %v px, height: %v px, refresh rate:%v Hz\n",
				j, vm.Width, vm.Height, vm.RefreshRate)

		}
	}

	primaryMonitor := pixelgl.PrimaryMonitor()

	primaryMonitorName := primaryMonitor.Name()
	fmt.Printf("\nprimary monitor name: %v\n", primaryMonitorName)

}

func main() {
	pixelgl.Run(run)
}
