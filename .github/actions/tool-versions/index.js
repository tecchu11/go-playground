const fs = require("fs");
const core = require("@actions/core");

async function run() {
  try {
    const filePath = core.getInput("file");
    const toolVersions = fs.readFileSync(filePath, "utf8").split("\n");

    toolVersions.forEach((line) => {
      if (line.trim()) {
        const [tool, version] = line.split(/\s+/);
        core.setOutput(tool, version);
      }
    });
  } catch (error) {
    core.setFailed(error.message);
  }
}

run();
