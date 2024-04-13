import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '1m', target: 20 },
        { duration: '5m', target: 50 },
        { duration: '1m', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<250'], // 95% of requests must complete below 250ms
    },
    rps: 2500,
};

const requestParams = [];

export default function () {
    if (requestParams.length === 0 || Math.random() < 0.05) {
        const params = {
            headers: {
                token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwNzI1OTIsImlzQWRtaW4iOnRydWV9.cAk-hqRR79HK2e6453lfOL5IjlP1WHkMPbKk9fFcBgQ',
                'Content-Type': 'application/json',
            },
        };
        const randomTagId = Math.floor(Math.random() * 100);
        const randomFeatureId = Math.floor(Math.random() * 100);
        const request = {
            tag_ids: [randomTagId],
            feature_id: randomFeatureId,
            content: {
                proident7: false,
                ad_fc4: -35552479,
                cupidatat_ad: 'labore consectetur in offi',
            },
            is_active: true,
        };

        const res = http.post(
            'http://localhost:8080/banner',
            JSON.stringify(request),
            params
        );
        requestParams.push({
            tag_id: randomTagId,
            feature_id: randomFeatureId,
        });
        check(res, { 'status is 201': (r) => r.status === 201 });
    } else {
        const params = {
            headers: {
                token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwNzI1OTIsImlzQWRtaW4iOnRydWV9.cAk-hqRR79HK2e6453lfOL5IjlP1WHkMPbKk9fFcBgQ',
            },
        };
        const randomIndex = Math.floor(Math.random() * requestParams.length);
        const res = http.get(
            `http://localhost:8080/user_banner?tag_id=${requestParams[randomIndex].tag_id}&feature_id=${requestParams[randomIndex].feature_id}&use_last_revision=false`,
            params
        );
        check(res, { 'status is 200': (r) => r.status === 200 });
    }
}