package lab.review;

import java.time.LocalDate;
import java.util.Map;

/** Minimal interval-repetition planner. Quality: 0=again, 1=hard, 2=good, 3=easy. */
public final class ReviewPlanner {
    private static final Map<Integer, Integer> DAYS = Map.of(0, 0, 1, 1, 2, 3, 3, 7);

    public static LocalDate nextReview(LocalDate today, int quality, int repetitions) {
        if (today == null) throw new IllegalArgumentException("today is required");
        if (!DAYS.containsKey(quality)) throw new IllegalArgumentException("quality must be 0..3");
        int baseDays = DAYS.get(quality);
        int adjusted = quality >= 2 ? Math.min(30, baseDays + Math.max(0, repetitions - 1) * 2) : baseDays;
        return today.plusDays(adjusted);
    }

    public static void main(String[] args) {
        int quality = Integer.parseInt(args.length > 0 ? args[0] : "2");
        int repetitions = Integer.parseInt(args.length > 1 ? args[1] : "1");
        System.out.println(nextReview(LocalDate.now(), quality, repetitions));
    }
}
