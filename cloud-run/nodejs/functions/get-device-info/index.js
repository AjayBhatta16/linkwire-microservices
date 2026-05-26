const functions = require("@google-cloud/functions-framework");
const deviceDetector = require("node-device-detector");

functions.http("handler", async (req, res) => {
    const data = parseRequest(req.body);

    if (data === null) {
        res.status(400).send("Bad Request: missing or invalid request data");
    }

    const { userAgent } = data;

    if (!userAgent) {
        res.status(400).send("Bad Request: missing userAgent field");
    }

    let detectedAgent = deviceDetector.detect(req.body.userAgent);

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

function parseRequest(body) {
    if (!body.message?.data) {
        return null;
    }
    
    try {
        const decoded = Buffer.from(body.message.data, "base64").toString("utf8");
        return JSON.parse(decoded);
    } 
    catch (err) {
        console.error(`Bad Request: could not decode/parse message data — ${err.message}`);
        return null;
    }
}