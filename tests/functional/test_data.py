import pytest
import json


# getResponse unwraps the data/error from json response.
# @expected shall be set to None only if
# the response result is just to generate a component for a test
# but not actually returning a test result.
def getResponse(responseText, expected=None):
    response = json.loads(responseText)
    if "error" in response:
        error = response["error"]
        if expected is None or \
                (expected is not None and error != expected["error"]):
            pytest.fail(f"Failed to run test.\nDetails: {error}")
        return None
    return response["data"]


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
          "data_to_store": [
              {
                  "created_at": "Sat Sep 23 2017 15:38:22.0123 GMT+0630",
                  "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132",
                  "data": {
                      "test1": "value1",
                      "test2": "value2"
                  }
              }
          ]
        },
        # Expected
        {
            "error": "",
            "data": "OK"
        }),
    (
        # Input data
        {
          "data_to_store": [
              {
                  "created_at": "Sat Sep 23 2017 15:38:23.0000 GMT+0630",
                  "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132",
                  "data": {
                      "test3": "value1",
                      "test4": "value2"
                  }
              },
              {
                  "created_at": "Sat Sep 23 2017 15:38:23.0123 GMT+0630",
                  "viewer_id": "863bda70-a0aa-45fc-bd0c-dedd81515292",
                  "data": {
                      "test5": "value1",
                      "test6": "value2"
                  }
              }
          ]
        },
        # Expected
        {
            "error": "",
            "data": "OK"
        }),
    (
        # Input data
        {
          "data_to_store": [
              {
                  "created_at": "Sat Sep 23 2017 15:38:23.0000 GMT+0630",
                  "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132",
                  "data": {
                      "test3": "value1",
                      "test4": "value2"
                  }
              }
          ]
        },
        # Expected
        {
            "error": "Failed to save data",
            "data": ""
        })
]

ids = ['Add single row', 'Add multiple rows', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddData(httpConnection, data, expected):
    try:
        r = httpConnection.POST("/add-data", data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return
    response = getResponse(r.text, expected)
    if response is None:
        return

    expectedData = expected["data"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expectedData}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "interval": "Sat Sep 23 2017 15:38:22.0000 GMT+0630",
            "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132"
        },
        # Expected
        {
            "error": "",
            "data": [
                {
                    "created_at": "2017-09-23T15:38:22.0123Z",
                    "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132",
                    "data": {
                        "test1": "value1",
                        "test2": "value2"
                    }
                },
                {
                    "created_at": "2017-09-23T15:38:23Z",
                    "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132",
                    "data": {
                        "test3": "value1",
                        "test4": "value2"
                    }
                }
            ]
        }),
    (
        # Input data
        {
            "interval": "Sat Sep 23 2017 15:38:24.0000 GMT+0630",
            "viewer_id": "4a8c1833-3e4d-4144-8aab-3133f2bdc132"
        },
        # Expected
        {
            "error": "",
            "data": []
        })
]

ids = ['Success', 'Empty']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetDataByViewer(httpConnection, data, expected):
    try:
        r = httpConnection.GET("/get-data-by-viewer", data)
    except Exception:
        pytest.fail("Failed to send GET request")
        return
    response = getResponse(r.text, expected)
    if response is None:
        return

    expectedData = expected["data"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expectedData}")
