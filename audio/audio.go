package audio

import "context"

// Функция для отмены контекста
var Cancel context.CancelFunc

// Контекст для начала и завершения передачи
var Context context.Context

// Ссылка на выбранное устройство захвата
var CurrentCaptureDevice *Device

// Все аудиоустройства
var AllDevices = make([]Device, 0)

// Выбранные устройства вывода
var CurrentReceiverDevices = make([]*Device, 0)

// Передаются ли в данный момент данные
var Streaming bool

// Ресэмплинг. Предполагается возможность алгоритма выбора ресэмплинга тут, но пока что TODO
func resample(input []byte, oldSampleRate, newSampleRate uint32, inputBitDepth int, outputBitDepth int) []byte {
	if oldSampleRate == newSampleRate && inputBitDepth == outputBitDepth {
		return input
	}
	// return resampleLerp(input, oldSampleRate, newSampleRate, inputBitDepth, outputBitDepth)
	return resampleTest(input, oldSampleRate, newSampleRate, inputBitDepth, outputBitDepth)
}

func resampleTest(input []byte, inputSampleRate, outputSampleRate uint32, inputBitDepth int, outputBitDepth int) []byte {

	var bitDepthRatio = float64(outputBitDepth) / float64(inputBitDepth)
	inputBytesPerSample := inputBitDepth / 8
	outputBytesPerSample := outputBitDepth / 8

	if (inputSampleRate == outputSampleRate && inputBitDepth != outputBitDepth) || (inputSampleRate > outputSampleRate && inputBitDepth < outputBitDepth) {
		output := make([]byte, int(float64(len(input))*bitDepthRatio))
		for i := 0; i < len(input); i += inputBytesPerSample {
			var sample int
			for j := 0; j < inputBytesPerSample; j++ {
				sample = sample | int(input[i+j])<<(8*j)
			}
			if inputBitDepth < outputBitDepth {
				sample = sample << (outputBitDepth - inputBitDepth)
				for j := 0; j < outputBytesPerSample; j++ {
					output[i*outputBytesPerSample/inputBytesPerSample+j] = byte(sample >> (8 * j))
				}
				inputBytesPerSample = outputBytesPerSample
			} else if inputBitDepth > outputBitDepth {
				sample = sample >> (inputBitDepth - outputBitDepth)
				for j := outputBytesPerSample - 1; j >= 0; j-- {
					output[i*outputBytesPerSample/inputBytesPerSample+j] = byte(sample >> (8 * j))
				}
			}
		}
		input = output
	}

	sampleRateRatio := float64(outputSampleRate) / float64(inputSampleRate)
	output := make([]byte, int(float64(len(input))*bitDepthRatio*sampleRateRatio))

	if inputSampleRate < outputSampleRate && inputBitDepth < outputBitDepth {
		for i := 0; i < len(output); i += outputBytesPerSample {
			idxD := float64(i/outputBytesPerSample) / sampleRateRatio
			a := int(idxD) * inputBytesPerSample
			b := a + inputBytesPerSample

			var sampleA int
			for j := 0; j < inputBytesPerSample; j++ {
				sampleA = sampleA | int(input[a+j])<<(8*j)
			}

			if b >= len(input) {
				sampleA = sampleA << 8
				for j := 0; j < outputBytesPerSample; j++ {
					output[i] = byte(sampleA >> 8 * j)
				}
				continue
			}

			// need to keep track of the sign
			sampleA = sampleToSignedInt(sampleA, inputBytesPerSample)

			var sampleB int
			for j := 0; j < inputBytesPerSample; j++ {
				sampleB = sampleB | int(input[b+j])<<(8*j)
			}

			sampleB = sampleToSignedInt(sampleB, inputBytesPerSample)

			bCF := idxD - float64(a/2)
			aCF := 1.0 - bCF
			// Linear interpolation
			resA := aCF*float64(sampleA) + bCF*float64(sampleB)

			resB := int(resA) << (outputBitDepth - inputBitDepth)
			for j := 0; j < outputBytesPerSample; j++ {
				output[i+j] = byte(resB >> (8 * j))
			}
		}
		return output
	}

	if inputSampleRate > outputSampleRate && inputBitDepth == outputBitDepth {
		for i := 0; i < len(output); i += outputBytesPerSample {
			idxD := float64(i/outputBytesPerSample) / sampleRateRatio
			a := int(idxD) * inputBytesPerSample
			b := a + inputBytesPerSample

			var sampleA int
			for j := 0; j < inputBytesPerSample; j++ {
				sampleA = sampleA | int(input[a+j])<<(8*j)
			}

			if b >= len(input) {
				sampleA = sampleA << 8
				for j := 0; j < outputBytesPerSample; j++ {
					output[i] = byte(sampleA >> 8 * j)
				}
				continue
			}

			sampleA = sampleToSignedInt(sampleA, inputBytesPerSample)

			var sampleB int
			for j := 0; j < inputBytesPerSample; j++ {
				sampleB = sampleB | int(input[b+j])<<(8*j)
			}

			sampleB = sampleToSignedInt(sampleB, inputBytesPerSample)

			bCF := idxD - float64(a/2)
			aCF := 1.0 - bCF

			resA := aCF*float64(sampleA) + bCF*float64(sampleB)

			resB := int(resA) >> (inputBitDepth - outputBitDepth)
			for j := outputBytesPerSample - 1; j >= 0; j-- {
				output[i+j] = byte(resB >> (8 * j))
			}
		}
	}
	return input
}

// Ресэмплинг с помощью линейной интерполяции
// Не работает в некоторых случаях, поэтому не используется
func resampleLerp(input []byte, inputSampleRate, outputSampleRate uint32, inputBitDepth int, outputBitDepth int) []byte {

	var bitDepthRatio = float64(outputBitDepth) / float64(inputBitDepth)
	inputBytesPerSample := inputBitDepth / 8
	outputBytesPerSample := outputBitDepth / 8

	var sampleRateRatio = float64(outputSampleRate) / float64(inputSampleRate)
	output := make([]byte, int(float64(len(input))*bitDepthRatio*sampleRateRatio))

	if inputSampleRate == outputSampleRate && inputBitDepth == outputBitDepth {
		copy(output, input)
		return output
	}

	if inputSampleRate == outputSampleRate && inputBitDepth != outputBitDepth {
		for i := 0; i < len(input); i += inputBytesPerSample {
			var sample int
			if inputBitDepth < outputBitDepth {
				for j := 0; j < inputBytesPerSample; j++ {
					sample = sample | int(input[i+j])<<(8*j)
				}

				sample = sample << (outputBitDepth - inputBitDepth)

				for j := 0; j < outputBytesPerSample; j++ {
					output[i*outputBytesPerSample/inputBytesPerSample+j] = byte(sample >> (8 * j))
				}

			} else if inputBitDepth > outputBitDepth {
				for j := 0; j < inputBytesPerSample; j++ {
					sample = sample | int(input[i+j])<<(8*j)
				}
				sample = sample >> (inputBitDepth - outputBitDepth)
				for j := outputBytesPerSample - 1; j >= 0; j-- {
					output[i*outputBytesPerSample/inputBytesPerSample+j] = byte(sample >> (8 * j))
				}
			}
		}
		return output
	}

	if inputSampleRate != outputSampleRate && inputBitDepth == outputBitDepth {
		for i := 0; i < len(output); i += outputBytesPerSample {
			idxD := float64(i/outputBytesPerSample) / sampleRateRatio
			a := int(idxD) * inputBytesPerSample
			b := a + inputBytesPerSample

			var sampleA int
			for j := 0; j < inputBytesPerSample; j++ {
				sampleA = sampleA | int(input[a+j])<<(8*j)
			}

			sampleA = sampleToSignedInt(sampleA, inputBytesPerSample)

			if b >= len(input) {
				for j := 0; j < outputBytesPerSample; j++ {
					output[i+j] = byte(sampleA >> (8 * j))
				}
				continue
			}

			var sampleB int
			for j := 0; j < inputBytesPerSample; j++ {
				sampleB = sampleB | int(input[b+j])<<(8*j)
			}
			sampleB = sampleToSignedInt(sampleB, inputBytesPerSample)

			bCF := idxD - float64(a/2)
			aCF := 1.0 - bCF
			// Linear interpolation
			resA := aCF*float64(sampleA) + bCF*float64(sampleB)
			resB := int(resA) >> (inputBitDepth - outputBitDepth)

			// Write to output
			for j := outputBytesPerSample - 1; j >= 0; j-- {
				output[i+j] = byte(int(resB) >> (8 * j))
			}
		}
		return output
	}

	if inputSampleRate > outputSampleRate && inputBitDepth < outputBitDepth {
		for i := 0; i < len(input); i += inputBytesPerSample {
			idxD := float64(i/inputBytesPerSample) * sampleRateRatio
			a := int(idxD) * outputBytesPerSample
			b := a + outputBytesPerSample

			var sample int
			for j := 0; j < inputBytesPerSample; j++ {
				sample = sample | int(input[i+j])<<(8*j)
			}
			sample = sample >> (inputBitDepth - outputBitDepth)

			if b >= len(output) {
				for j := outputBytesPerSample - 1; j >= 0; j-- {
					output[a+j] = byte(sample >> (8 * j))
				}
				continue
			}

			var nextSample int
			for j := 0; j < outputBytesPerSample; j++ {
				nextSample = nextSample | int(output[b+j])<<(8*j)
			}

			bCF := idxD - float64(a/outputBytesPerSample)
			aCF := 1.0 - bCF

			resA := aCF*float64(sample) + bCF*float64(nextSample)

			for j := outputBytesPerSample - 1; j >= 0; j-- {
				output[a+j] = byte(int(resA) >> (8 * j))
			}
		}
		return output
	}

	if inputSampleRate < outputSampleRate && inputBitDepth > outputBitDepth {
		for i := 0; i < len(output); i += inputBytesPerSample {
			idxD := float64(i/outputBytesPerSample) / sampleRateRatio
			a := int(idxD) * 3
			b := a + 3

			var sampleA int
			sampleA = sampleA | int(input[a]) | int(input[a+1])<<8 | int(input[a+2])<<16
			sampleA = sampleToSignedInt(sampleA, 3)

			if b >= len(input) {

				sampleA = sampleToSignedInt(sampleA, 3)
				output[i] = byte(sampleA)
				output[i+1] = byte(sampleA >> 8)
				continue
			}

			var sampleB int
			sampleB = sampleB | int(input[b]) | int(input[b+1])<<8 | int(input[b+2])<<16
			sampleB = sampleToSignedInt(sampleB, 3)

			bCF := idxD - float64(a/3)
			aCF := 1.0 - bCF

			resA := aCF*float64(sampleA) + bCF*float64(sampleB)

			resB := sampleToSignedInt(int(resA), 3) >> 8

			output[i] = byte(resB >> 16)
			output[i+1] = byte(resB >> 8)
		}
		return output
	}
	return output
}

// Перевод из int8, 16, 24, 32 с учётом знака
func sampleToSignedInt(sample int, bytesPerSample int) int {

	switch bytesPerSample {
	case 1:
		sample = int(int8(sample))
	case 2:
		sample = int(int16(sample))
	case 3:
		// For 24-bit audio, the sign bit is in the third byte
		// if (sample & int(0x00800000)) != 0 {
		// 	sample |= int(uint(0xFF000000)) // Extend sign bit to 32 bits
		// }
		sample = int(int32(sample))
	case 4:
		sample = int(int32(sample))
	}
	return sample
}
