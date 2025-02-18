# Text to Speech with FFmpeg and gTTS Py Integration

This project integrates FFmpeg and gTTS Py with a Go application to process text to speech.

## Prerequisites

1. Install Docker on your OS or server. Follow the official Docker installation guide: [Docker Documentation](https://docs.docker.com/get-docker/).

## Steps to Run

1. **Build the Docker Image**  
   Run the following command in the project directory:
   ```bash
   docker build -t generatetts .
   ```

2. **Prepare Input Files**  
   Ensure the input files are present in the same repository as the application.

3. **Run the Application**  
   Use the following command to run the Go application:
   ```bash
   docker run --rm -v "$(pwd):/output" tts-converter
   ```

## Notes

- The application expects audio files for processing to be available in the project directory.
- Docker is used for building the environment with FFmpeg integration.

## License

MIT

