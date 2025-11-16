import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomItem, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
    vus: 200,
    duration: '20s',
};

function generateRandomString(prefix = '', length = 6) {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = prefix;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return result;
}

function generateRandomUser() {
    const randomUserId = generateRandomString('u', 5);
    const randomUsername = generateRandomString('', 8);
    return {
        user_id: randomUserId,
        username: randomUsername,
        is_active: true
    };
}

function generateRandomTeam(numMembers = 2) {
    const randomTeamName = generateRandomString('', 10); 
    const members = [];
    for (let i = 0; i < numMembers; i++) {
        members.push(generateRandomUser());
    }
    return {
        team_name: randomTeamName,
        members: members
    };
}

const baseUrl = 'http://localhost:8080/api/v1'; 

export default function () {
    const team = generateRandomTeam(randomIntBetween(2, 10))
    let res = http.post(`${baseUrl}/team/add`, JSON.stringify(team), {
        headers: { 'Content-Type': 'application/json'},
    });
    check(res, { 'team add status 201': (r) => r.status === 201 || r.status === 400 });

    sleep(randomIntBetween(1, 3));
}
