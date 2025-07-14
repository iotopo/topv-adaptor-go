#!/bin/bash

echo "Testing TopV Adaptor Go API..."

echo
echo "1. Testing query devices..."
curl -X GET http://localhost:8080/api/query_devices \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1"}'

echo
echo "2. Testing query points..."
curl -X GET http://localhost:8080/api/query_points \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1","parentTag":"group1.dev1"}'

echo
echo "3. Testing find last..."
curl -X GET http://localhost:8080/api/find_last \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1","tag":"group1.dev1.a","device":false}'

echo
echo "4. Testing query history..."
curl -X POST http://localhost:8080/api/query_history \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1","tag":["group1.dev1.a"],"start":"2024-01-01T00:00:00Z","end":"2024-01-01T23:59:59Z"}'

echo
echo "5. Testing set value..."
curl -X POST http://localhost:8080/api/set_value \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1","tag":"group1.dev1.a","value":"25.5","time":1640995200000}'

echo
echo "API tests completed!" 