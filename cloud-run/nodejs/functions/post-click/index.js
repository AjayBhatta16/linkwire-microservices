const functions = require("@google-cloud/functions-framework");
const http = require("http");
const deviceDetector = require("node-device-detector");

const { Firestore } = require('@google-cloud/firestore');

functions.http("handler", async (req, res) => {
    const data = parseRequest(req.body);

    if (data === null) {
        res.status(400).send("Bad Request: missing or invalid message data");
    }

    const { ipAddress, userAgent, linkID } = data;

    const deviceInfo = extractDeviceInfo(userAgent);
    const ipInfo = await extractIPInfo(ipAddress);

    const click = {
        clickID: `${linkID}-${Date.now()}`,
        linkID: linkID,
        ip: ipAddress,
        timestamp: Date.now(),
        userAgent: userAgent,
        os: `${device.os.name} ${device.os.version}`,
        client: `${device.client.type} - ${device.client.name} ${device.client.version}`,
        device: `${device.device.type} - ${device.device.type} ${device.device.model}`,
        location: `${ipInfo.city}, ${ipInfo.regionName}, ${ipInfo.country}`,
        isp: ipInfo.isp,
        asn: ipInfo.as,
        mobile: ipInfo.mobile,
        proxy: ipInfo.proxy,
        hosting: ipInfo.hosting,
    };

    await saveClick(click);

    res.status(200).send("Click processed successfully");
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

function extractDeviceInfo(userAgent) {
    const detector = new deviceDetector();
    return detector.detect(userAgent);
}

async function extractIPInfo(ipAddress) {
    const url = `http://ip-api.com/json/${encodeURIComponent(ipAddress)}`
        + '?fields=city,regionName,country,isp,as,mobile,proxy,hosting';

    return new Promise((resolve, reject) => {
        http.get(url, (res) => {
            let data = "";

            res.on("data", (chunk) => {
                data += chunk;
            });

            res.on("end", () => {
                if (res.statusCode !== 200) {
                    return reject(new Error(`ip-api.com returned status ${res.statusCode}`));
                }

                try {
                    resolve(JSON.parse(data));
                } catch (err) {
                    reject(new Error(`Failed to parse response from ip-api.com: ${err.message}`));
                }
            });
        }).on("error", (err) => {
            reject(new Error(`HTTP request failed: ${err.message}`));
        });
    });
}

async function saveClick(click) {
    const client = new Firestore();
    const docRef = client.collection("clicks").doc(click.clickID);
    await docRef.set(click);
}