package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/yesfeedme-archive/go-effects/pkg/effects"
)

func main() {
	effect := flag.String("effect", "", "The name of the effect to apply. Values are 'oil|sobel|gaussian|cartoon|pixelate'")
	flag.Parse()
	validateFlags(*effect)

	var inPath, outPath string
	inPath = flag.Arg(0)
	outPath = flag.Arg(1)

	img, err := effects.LoadImage(inPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outImg := runEffect(img, *effect)
	err = outImg.Save(outPath, effects.SaveOpts{ClipToBounds: true})
	if err != nil {
		fmt.Println("Failed to save modified image:", err)
		os.Exit(1)
	}

}

func runGaussian(img *effects.Image) *effects.Image {
	kernelSize, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("invalid kernelSize value")
		os.Exit(1)
	}
	sigma, err := strconv.ParseFloat(flag.Arg(3), 64)
	if err != nil {
		fmt.Println("invalid sigma value")
		os.Exit(1)
	}
	gaussian := effects.NewGaussian(kernelSize, sigma)
	outImg, err := gaussian.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runSobel(img *effects.Image) *effects.Image {
	threshold, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("invalid threshold value")
		os.Exit(1)
	}
	invert, err := strconv.ParseBool(flag.Arg(3))
	if err != nil {
		fmt.Println("invalid invert value")
		os.Exit(1)
	}
	sobel := effects.NewSobel(threshold, invert)
	outImg, err := sobel.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runPencil(img *effects.Image) *effects.Image {
	blurFactor, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("Invalid blurFactor value:", err)
		os.Exit(1)
	}

	pencil := effects.NewPencil(blurFactor)
	outImg, err := pencil.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runBrightness(img *effects.Image) *effects.Image {
	offset, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("Invalid offset value:", err)
		os.Exit(1)
	}

	brightness := effects.NewBrightness(offset)
	outImg, err := brightness.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runOil(img *effects.Image) *effects.Image {
	filterSize, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("Invalid filterSize value:", err)
		os.Exit(1)
	}
	if filterSize <= 3 {
		fmt.Println("FilterSize must be at least 3")
		os.Exit(1)
	}
	levels, err := strconv.Atoi(flag.Arg(3))
	if err != nil {
		fmt.Println("Invalid levels value:", err)
		os.Exit(1)
	}
	if levels < 1 {
		fmt.Println("Levels must be at least 1")
		os.Exit(1)
	}

	oil := effects.NewOilPainting(filterSize, levels)
	outImg, err := oil.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runCartoon(img *effects.Image) *effects.Image {
	blurStrength, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("Invalid blurStrength value")
		os.Exit(1)
	}
	edgeThreshold, err := strconv.Atoi(flag.Arg(3))
	if err != nil {
		fmt.Println("Invalid edgeThreshold value")
		os.Exit(1)
	}
	oilFilterSize, err := strconv.Atoi(flag.Arg(4))
	if err != nil {
		fmt.Println("Invalid oilFilterSize value")
		os.Exit(1)
	}
	oilLevels, err := strconv.Atoi(flag.Arg(5))
	if err != nil {
		fmt.Println("Invalid oilLevels value")
		os.Exit(1)
	}
	opts := effects.CTOpts{
		BlurKernelSize: blurStrength,
		EdgeThreshold:  edgeThreshold,
		OilFilterSize:  oilFilterSize,
		OilLevels:      oilLevels,
		DebugPath:      "",
	}
	cartoon := effects.NewCartoon(opts)
	outImg, err := cartoon.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runPixelate(img *effects.Image) *effects.Image {
	blockSize, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("Invalid blockSize value")
		os.Exit(1)
	}
	pixelate := effects.NewPixelate(blockSize)
	outImg, err := pixelate.Apply(img, 0)
	if err != nil {
		fmt.Println("Failed to apply effect:", err)
		os.Exit(1)
	}
	return outImg
}

func runEffect(img *effects.Image, effect string) *effects.Image {
	switch effect {
	case "brightness":
		return runBrightness(img)
	case "cartoon":
		return runCartoon(img)
	case "gaussian":
		return runGaussian(img)
	case "oil":
		return runOil(img)
	case "pencil":
		return runPencil(img)
	case "pixelate":
		return runPixelate(img)
	case "sobel":
		return runSobel(img)
	}
	return nil
}

func validateFlags(effect string) {
	switch effect {
	case "brightness":
		if len(flag.Args()) != 3 {
			fmt.Println("The brightness effect requires 3 args, input path, output path offset")
			fmt.Println("Sample usage: goeffects -effect=brightness mypic.jpg mypic-lighten.jpg 200")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "cartoon":
		if len(flag.Args()) != 6 {
			fmt.Println("The cartoon effect requires 6 args, input path, output path, blurStrength, edgeThreshold, oilBoldness, oilLevels")
			fmt.Println("Sample usage: goeffects -effect=cartoon mypic.jpg mypic-cartoon.jpg 21 40 15 15")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "gaussian":
		if len(flag.Args()) != 4 {
			fmt.Println("The gaussian effect requires 4 args, input path, output path, kernelSize, sigma")
			fmt.Println("Sample usage: goeffects -effect=gaussian mypic.jpg mypic-gaussian.jpg 9 1")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "oil":
		if len(flag.Args()) != 4 {
			fmt.Println("The oil effect requires 4 args, input path, output path, filterSize, levels")
			fmt.Println("Sample usage: goeffects -effect=oil mypic.jpg mypic-oil.jpg 5 30")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "pencil":
		if len(flag.Args()) != 3 {
			fmt.Println("The pencil effect requires 3 args, input path, output path blurFactor")
			fmt.Println("Sample usage: goeffects -effect=pencil mypic.jpg mypic-pencil.jpg")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "pixelate":
		if len(flag.Args()) != 3 {
			fmt.Println("The pixelate effect requires 3 args, input path, output path, block size")
			fmt.Println("Sample usage: goeffects -effect=pixelate mypic.jpg mypic-pixelate.jpg 12")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "sobel":
		if len(flag.Args()) != 4 {
			fmt.Println("The sobel effect requires 4 args, input path, output path, threshold invert")
			fmt.Println("Sample usage: goeffects -effect=sobel mypic.jpg mypic-sobel.jpg 100 false")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "":
		fmt.Println("The effect option is required")
		flag.PrintDefaults()
		os.Exit(1)

	default:
		fmt.Println("Unknown effect option value")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
