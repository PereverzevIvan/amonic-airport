# async def test_basic(api):
#     data = {"url": "http://example.com?id=123"}
#     response = await api.post('/v1/make-shorter', json=data)
#     assert response.status == 200

#     response_json = response.json()
#     assert 'short_url' in response_json

#     request_url = response_json['short_url']
#     request_url = request_url[request_url.rfind('/'):]

#     response = await api.get(request_url)
#     assert response.status == 200


# async def test_bad_request(api):
#     data = {}
#     response = await api.post('/v1/make-shorter', json=data)
#     assert response.status == 400


# async def test_redirect_not_found(api):
#     response = await api.get('/unknown-id')
#     assert response.status == 404
