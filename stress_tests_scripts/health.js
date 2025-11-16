import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

export let successRate = new Rate("successful_requests");

export let options = {
    vus: 5,
    duration: '60s',
};


export default function () {
    const res = http.get(`http://localhost:8080/health`);
    const result = check(res, {
        "status is 200": (r) => r.status === 200,
    });

    successRate.add(result);

    sleep(1);
}
