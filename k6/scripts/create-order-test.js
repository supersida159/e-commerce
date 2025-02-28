import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  // Configure test scenarios
  scenarios: {
    ramping_vus_scenario: {
      executor: 'ramping-vus',
      startVUs: 5, // Start with 5 VUs
      stages: [
        { duration: '10s', target: 10 }, // Ramp up to 10 VUs over 10 seconds
        { duration: '20s', target: 10 }, // Stay at 10 VUs for 20 seconds
        { duration: '20s', target: 40 }, // Ramp up to 40 VUs over 20 seconds
        { duration: '10s', target: 0 },  // Ramp down to 0 VUs over 10 seconds
      ],
      gracefulRampDown: '0s', // Wait time for iterations to finish before stopping VUs
    },
  },
  thresholds: {
    // Set performance thresholds
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.01'],   // Less than 1% of requests should fail
  },
};

// JWT token - storing this directly in the script for simplicity
// In production, consider more secure options for handling tokens
const jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7InVzZXJfaWQiOjEsInJvbGUiOiJ1c2VyIn0sImV4cCI6MTc0MDczNDc1NywiaWF0IjoxNzQwNDgyNzU3fQ.ZItCTGdY-0btEy1BxeZAZfejiKJSWWI0dNgLq7Niuog';

export default function() {
  // Define the API endpoint
  // Note: When running in Docker, use the service name instead of localhost
  const url = 'http://api-services:8090/api/v1/order/Private/createOrder';
  
  // Prepare request headers with JWT token
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwtToken}`
    },
  };

  // Prepare order payload
  // Generate some random values to avoid duplicate orders
  const orderId = Math.floor(Math.random() * 1000000);
  const payload = JSON.stringify({
    "customer_name": `Test User ${orderId}`,
    "customer_phone": "+1234567890",
    "notes": "K6 load test order",
    "address_id": 1,
    "shipping": {
      "method": "Standard",
      "cost": 5.00,
      "estimated_delivery": "2025-03-01"
    }
  });

  // Send POST request to create order
  const response = http.post(url, payload, params);

  // Verify the response
  check(response, {
    'is status 200 or 202': (r) => r.status === 200 || r.status === 202,
    'has order ID': (r) => JSON.parse(r.body).order_id !== undefined,
    'response time < 200ms': (r) => r.timings.duration < 200
  });

  // Log details for debugging (optional)
  if (response.status !== 200 && response.status !== 202) {
    console.log(`Error: ${response.status} - ${response.body}`);
  }

  // Add a small pause between requests
  sleep(1);
}



// k6 cmd loadtest
// docker-compose run k6 run -o influxdb=http://influxdb:8086/k6 /scripts/create-order-test.js

