import http from 'k6/http';
import { sleep } from 'k6';
import { SharedArray }  from "k6/data";
import { scenario } from 'k6/execution';

const data = new SharedArray(
    "data",
    function () {
        let ids = []
        for (var i = 0; i < 1000; i++) {
            ids.push(i)
        }
        return ids
    }
)

export const options = {
    scenarios: {
        use_data: {
            executor: 'shared-iterations',
            vus: 100,
            iterations: data.length,
            maxDuration: '10s',
        }
    },
    thresholds: {
        'http_req_duration': ['p(99)<1000'], // 99% of requests must complete below 1.0s
    }
}

export default function () {
    const user = data[scenario.iterationInTest];
    http.get('http://localhost:8081/redis1?key=' + user);
    sleep(1);
}