module.exports = {
    * beforeSendRequest(requestDetail) {
        console.log(requestDetail.url);
        if (requestDetail.url.indexOf('/2016-08-15/services') !== -1) {
            const localResponse = {
                statusCode: 200,
                header: {
                    'Content-Type': 'application/json'
                },
                body: `{
                "services": [{
                  "serviceName": "demo",
                  "description": ""
                }]
              }`
            };
            return {
                response: localResponse
            };
        }
    },
};