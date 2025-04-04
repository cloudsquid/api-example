const fs = require("fs");
const path = require("path");
const axios = require("axios");
const base64 = require("base64-js");
require("dotenv").config();

class CloudsquidClient {
  constructor(apikey, endpoint, sourceId) {
    this.apikey = apikey;
    this.endpoint = endpoint;
    this.sourceId = sourceId;
    this.client = axios.create({
      baseURL: endpoint,
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": apikey,
      },
    });
  }

  async uploadFile({ mimetype, filename, filetype, file }) {
    console.log("Uploading file");
    const url = path.join("datasources", this.sourceId, "documents");
    const payload = { mimetype, filename, file_type: filetype, file };

    const response = await this.client.post(url, payload);
    console.log("Successfully sent out upload request");
    return response.data;
  }

  async runFile({ fileId, pipeline }) {
    console.log("Running file");
    const url = path.join("datasources", this.sourceId, "run");
    const payload = { file_id: fileId, pipeline };

    const response = await this.client.post(url, payload);
    console.log("Successfully sent out request to run file");
    return response.data;
  }

  async getStatus(runId) {
    console.log("Getting status");
    const url = path.join("datasources", this.sourceId, "run", runId);

    const response = await this.client.get(url);
    console.log("Successfully got status of file");
    return response.data;
  }
}

async function main() {
  const config = {
    CsKey: process.env.CLOUDSQUID_API_KEY,
    CsEndpoint: process.env.CLOUDSQUID_API_ENDPOINT,
    CsSourceID: process.env.CLOUDSQUID_AGENT_ID,
  };

  console.log(config);

  const filePath = process.argv[2];
  if (!filePath) {
    console.error("Usage: node index.js <file-path>");
    process.exit(1);
  }

  const client = new CloudsquidClient(
    config.CsKey,
    config.CsEndpoint,
    config.CsSourceID
  );

  const fileContent = fs.readFileSync(filePath);
  const fileName = path.basename(filePath);
  const encodedFile = base64.fromByteArray(fileContent);

  const uploadPayload = {
    mimetype: "application/pdf",
    filename: fileName,
    filetype: "binary",
    file: encodedFile,
  };

  const uploadResponse = await client.uploadFile(uploadPayload);
  console.log("Upload Response:", uploadResponse);

  const runPayload = {
    fileId: uploadResponse.file_id,
    pipeline: "cloudsquid-flash",
  };

  const runResponse = await client.runFile(runPayload);
  console.log("Run Response:", runResponse);

  let extraction;
  while (true) {
    const statusResponse = await client.getStatus(runResponse.run_id);
    console.log("Status Response:", statusResponse);

    if (statusResponse.status === "done") {
      extraction = JSON.stringify(statusResponse.result, null, 2);
      break;
    }

    if (statusResponse.status === "error") {
      console.error("Error in processing:", statusResponse.result);
      process.exit(1);
    }

    await new Promise((resolve) => setTimeout(resolve, 2000)); // Delay to avoid too many requests
  }

  console.log("Final result:\n", extraction);
}

main().catch((err) => {
  console.error("Error:", err);
  process.exit(1);
});
