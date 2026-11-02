.PHONY: web api audio java order

web:
	cd apps/web && node server.mjs

api: order
	cd services/api && PHRASE_ORDER_BIN="../../engines/order/target/release/phrase-order" go run .

audio:
	cmake -S engines/audio -B engines/audio/build
	cmake --build engines/audio/build
	ctest --test-dir engines/audio/build --output-on-failure

java:
	cd services/review && mkdir -p build/classes && javac -d build/classes ReviewPlanner.java ReviewPlannerTest.java
	cd services/review && java -cp build/classes lab.review.ReviewPlannerTest

order:
	cd engines/order && cargo test && cargo build --release
