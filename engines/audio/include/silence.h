#ifndef PHRASE_LAB_SILENCE_H
#define PHRASE_LAB_SILENCE_H
#include <stddef.h>
#include <stdint.h>
#ifdef __cplusplus
extern "C" {
#endif
size_t count_quiet_frames(const int16_t *samples, size_t sample_count, size_t frame_size, int16_t threshold);
#ifdef __cplusplus
}
#endif
#endif
