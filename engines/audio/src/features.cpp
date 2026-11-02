#include "silence.h"
#include <cmath>
#include <cstdint>
#include <iostream>
#include <vector>

int main() {
    // Prototype input is a short PCM sample sequence; WAV decoding is a future adapter.
    std::vector<int16_t> samples;
    int value;
    while (std::cin >> value) {
        if (value < -32768 || value > 32767) {
            std::cerr << "sample out of int16 range\n";
            return 2;
        }
        samples.push_back(static_cast<int16_t>(value));
    }
    if (samples.empty()) {
        std::cerr << "provide signed 16-bit PCM samples on stdin\n";
        return 2;
    }
    long double energy = 0;
    for (int16_t sample : samples) energy += static_cast<long double>(sample) * sample;
    const double rms = std::sqrt(static_cast<double>(energy / samples.size()));
    const size_t quiet = count_quiet_frames(samples.data(), samples.size(), 160, 250);
    std::cout << "{\"sampleCount\":" << samples.size() << ",\"rms\":" << rms
              << ",\"quietFrames\":" << quiet << "}\n";
}
