#include "silence.h"

size_t count_quiet_frames(const int16_t *samples, size_t sample_count, size_t frame_size, int16_t threshold) {
    if (samples == NULL || frame_size == 0) return 0;
    size_t quiet = 0;
    for (size_t start = 0; start + frame_size <= sample_count; start += frame_size) {
        long long energy = 0;
        for (size_t i = 0; i < frame_size; ++i) {
            long long value = samples[start + i];
            energy += value * value;
        }
        long long mean_square = energy / (long long)frame_size;
        long long threshold_square = (long long)threshold * threshold;
        if (mean_square <= threshold_square) ++quiet;
    }
    return quiet;
}
