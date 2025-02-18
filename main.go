package main

import (
	"fmt"
	"os"
	"os/exec"
)

func textToSpeech(textInput string, lang string, outputFile string, gender string) error {
	outputDir := "/output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	}

	outputPath := fmt.Sprintf("%s/%s", outputDir, outputFile)
	tempFile := fmt.Sprintf("%s/temp.mp3", outputDir)

	// Gunakan gTTS (tetap dengan suara default)
	pythonCmd := exec.Command("python3", "-c", fmt.Sprintf(`
from gtts import gTTS
tts = gTTS("%s", lang="%s")
tts.save("%s")
`, textInput, lang, tempFile))

	pythonCmd.Stderr = os.Stderr
	pythonCmd.Stdout = os.Stdout

	fmt.Println("Running gTTS...")
	if err := pythonCmd.Run(); err != nil {
		return fmt.Errorf("failed to run gTTS: %v", err)
	}

	// Jika gender = "pria", ubah pitch suara dengan ffmpeg
	if gender == "male" {
		fmt.Println("Modifying voice for male effect...")
		modifiedFile := fmt.Sprintf("%s/temp_male.mp3", outputDir)
		ffmpegPitchCmd := exec.Command("ffmpeg", "-i", tempFile, "-af", "asetrate=18000,atempo=1.5", modifiedFile)
		ffmpegPitchCmd.Stderr = os.Stderr
		ffmpegPitchCmd.Stdout = os.Stdout
		if err := ffmpegPitchCmd.Run(); err != nil {
			return fmt.Errorf("failed to modify pitch: %v", err)
		}
		// Gunakan suara yang sudah dimodifikasi
		tempFile = modifiedFile
	}

	// Konversi ke AMR dengan ffmpeg
	ffmpegCmd := exec.Command("ffmpeg", "-i", tempFile, "-ar", "8000", "-ac", "1", "-ab", "12.2k", "-c:a", "libopencore_amrnb", outputPath)
	ffmpegCmd.Stderr = os.Stderr
	ffmpegCmd.Stdout = os.Stdout

	fmt.Println("Running ffmpeg conversion to AMR...")
	if err := ffmpegCmd.Run(); err != nil {
		return fmt.Errorf("failed to convert to AMR: %v", err)
	}

	// Hapus file sementara
	os.Remove(tempFile)

	fmt.Println("Conversion completed successfully! File saved at:", outputPath)
	return nil
}

func main() {
	text := "Halo, ini adalah contoh teks ke suara menggunakan gTTS."
	outputFile := "output.amr"
	lang := "zh"
	gender := "male" // Bisa "pria" atau "wanita"

	if err := textToSpeech(text, lang, outputFile, gender); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
