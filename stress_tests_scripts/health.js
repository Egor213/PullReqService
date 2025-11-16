import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

export let successRate = new Rate("successful_requests");

export let options = {
    vus: 300,
    duration: '60s',
};


export default function (data) {
    const teamName = randomTeam();

    const res = http.get(`http://localhost:8080/health}`, {
        headers: {
            Authorization: `Bearer ${data.token}`,
        },
    });

    const result = check(res, {
        "status is 200": (r) => r.status === 200,
        "response has correct team_name": (r) => r.json("team_name") === teamName,
    });

    successRate.add(result);

    sleep(1);
}
