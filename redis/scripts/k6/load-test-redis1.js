import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
    stages: [
        {duration: '5m', target: 100},
        {duration: '10m', target: 100},
    ],
    thresholds: {
        'http_req_duration': ['p(99)<1000'], // 99% of requests must complete below 1.0s
    }
}

export default function () {
    http.get('http://localhost:8081/redis1');
    sleep(1);
}