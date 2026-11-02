package lab.review;

import java.time.LocalDate;

public final class ReviewPlannerTest {
    public static void main(String[] args) {
        LocalDate today = LocalDate.of(2026, 11, 2);
        if (!ReviewPlanner.nextReview(today, 0, 1).equals(today)) throw new AssertionError("again should be today");
        if (!ReviewPlanner.nextReview(today, 2, 1).equals(today.plusDays(3))) throw new AssertionError("good should be 3 days");
        if (!ReviewPlanner.nextReview(today, 3, 2).equals(today.plusDays(9))) throw new AssertionError("easy interval");
    }
}
