import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

export let successRate = new Rate("successful_requests");

export let options = {
    vus: 200,
    duration: '10s',
};

const teams = [
    "backend", "frontend", "mobile", "payments", "analytics",
    "core", "devops", "qa", "design", "ml", "ds", "hrtools",
    "search", "maps", "ads", "marketplace", "billing", "infra",
    "crypto", "api"
];

function randomTeam() {
    return teams[Math.floor(Math.random() * teams.length)];
}

export function setup() {
    const payload = JSON.stringify({
        user_id: "u1",
        role: "admin"
    });

    const headers = { "Content-Type": "application/json" };

    const res = http.post("http://localhost:8080/auth/login", payload, { headers });
    check(res, {
        "login status is 200": (r) => r.status === 200,
        "has token": (r) => r.json("access_token") !== undefined,
    });

    const token = res.json("access_token");
    return { token };
}

export default function (data) {
    const teamName = randomTeam();

    const res = http.get(`http://localhost:8080/api/v1/team/get?team_name=${teamName}`, {
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
