import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

const publicViewTrend = new Trend('public_view_duration');
const bookingSuccessRate = new Rate('booking_success_rate');
const bookingConflictRate = new Rate('booking_conflict_rate');

export const options = {
    stages: [
        { duration: '10s', target: 5 },
        { duration: '20s', target: 30 },
        { duration: '10s', target: 0 },
    ],
    thresholds: {
        'public_view_duration': ['p(95)<500'],
        'booking_success_rate': ['rate>0.01'],
        'http_req_failed': ['rate<0.01'],
    },
};

export function setup() {
    const userEmail = `load_${Date.now()}@test.com`;
    const password = 'loadtest123';

    let res = http.post(`${BASE_URL}/auth/register`, JSON.stringify({
        email: userEmail,
        password: password,
    }), { headers: { 'Content-Type': 'application/json' } });

    if (res.status !== 201) {
        throw new Error(`Registration failed: ${res.status}`);
    }

    res = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
        email: userEmail,
        password: password,
    }), { headers: { 'Content-Type': 'application/json' } });

    if (res.status !== 200) {
        throw new Error(`Login failed: ${res.status}`);
    }

    const token = res.json('token');

    const wishlists = [];
    for (let i = 0; i < 3; i++) {
        const wlRes = http.post(`${BASE_URL}/wishlists`, JSON.stringify({
            title: `Load Test List ${i}`,
            description: `Generated for k6 testing ${Date.now()}`,
            event_date: new Date(Date.now() + 86400000 * (7 + i)).toISOString().split('T')[0],
        }), {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        if (wlRes.status === 201) {
            const wl = wlRes.json();
            for (let j = 0; j < 2; j++) {
                const itemRes = http.post(`${BASE_URL}/wishlists/${wl.id}/items`, JSON.stringify({
                    title: `Gift ${i}-${j}`,
                    description: `Description ${i}-${j}`,
                    product_link: `https://example.com/${i}/${j}`,
                    priority: (j % 5) + 1,
                }), {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    }
                });
                if (itemRes.status === 201) {
                    const item = itemRes.json();
                    wl.items = wl.items || [];
                    wl.items.push(item);
                }
            }
            wishlists.push(wl);
        }
    }

    return {
        wishlists: wishlists.map(w => ({
            token: w.access_token,
            items: w.items || [],
        })),
    };
}

export default function (data) {
    if (!data.wishlists || data.wishlists.length === 0) {
        console.warn('No wishlists available for testing');
        return;
    }

    const wl = data.wishlists[Math.floor(Math.random() * data.wishlists.length)];

    const viewStart = Date.now();
    const viewRes = http.get(`${BASE_URL}/public/wishlists/${wl.token}`);
    publicViewTrend.add(Date.now() - viewStart);

    check(viewRes, {
        'public view status is 200': (r) => r.status === 200,
    });

    if (wl.items && wl.items.length > 0) {
        const item = wl.items[Math.floor(Math.random() * wl.items.length)];
        const bookRes = http.post(`${BASE_URL}/public/wishlists/${wl.token}/items/${item.id}/book`);

        const booked = bookRes.status === 204;
        const conflict = bookRes.status === 409;

        bookingSuccessRate.add(booked);
        bookingConflictRate.add(conflict);

        check(bookRes, {
            'booking response is 204 or 409': (r) => r.status === 204 || r.status === 409,
        });
    }

    sleep(Math.random() * 2 + 0.5);
}

export function teardown(data) {
    console.log(`Test completed. Created ${data?.wishlists?.length || 0} wishlists.`);
}