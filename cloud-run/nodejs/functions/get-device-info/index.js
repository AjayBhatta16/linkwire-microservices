const functions = require("@google-cloud/functions-framework");
const deviceDetector = require("node-device-detector");

functions.http("handler", async (req, res) => {
    if (req.body === null) {
        res.status(400).send("Bad Request: missing or invalid request data");
    }

    const { userAgent } = req.body;

    if (!userAgent) {
        res.status(400).send("Bad Request: missing userAgent field");
    }

    let detectedAgent = deviceDetector.detect(userAgent);

    res.status(200).json({
        type: detectedAgent.device.type,
        brand: detectedAgent.device.brand,
        model: detectedAgent.device.model,
        osName: detectedAgent.os.name,
        osVersion: detectedAgent.os.version,
        clientType: detectedAgent.client.type,
        clientName: detectedAgent.client.name,
        clientVersion: detectedAgent.client.version,
    });
});