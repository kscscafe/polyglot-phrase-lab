#include "silence.h"
#include <cassert>
#include <cstdint>

int main() {
    const int16_t data[] = {0, 0, 0, 0, 1000, 1000, 1000, 1000};
    assert(count_quiet_frames(data, 8, 4, 10) == 1);
    assert(count_quiet_frames(data, 8, 0, 10) == 0);
    assert(count_quiet_frames(nullptr, 8, 4, 10) == 0);
}
