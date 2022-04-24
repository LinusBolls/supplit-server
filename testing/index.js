const { spawn } = require("child_process");
const { join } = require("path");

const axios = require("axios");

// const server = spawn("go", ["run", "main.go"], {
//   env: { ...process.env, GO111MODULE: "off" },
// });

// server.stdout.on("data", (data) => {
//   console.log(`stdout: ${data}`);

//   if (data.toString().startsWith("Listening on")) {
//     console.log("starting tests");
//     main();
//   }
// });

// server.stderr.on("data", (data) => {
//   console.error(`stderr: ${data}`);
// });

// server.on("close", (code) => {
//   console.log(`child process exited with code ${code}`);
// });

const {
  promises: { readFile },
} = require("fs");

const schema = {
  in: { columns: ["price", "discount"] },
  out: { columns: ["newPrice"] },
  nodes: ["multiply"],
  noodles: [
    [
      [0, 0],
      [2, 0],
    ],
    [
      [1, 0],
      [2, 1],
    ],
    [
      [2, 2],
      [3, 0],
    ],
  ],
};

async function main() {
  const file = await readFile(join(__dirname, "./data.csv"));

  console.log(file.toString());

  const dater = {
    csv: file.toString(),
    schema,
  };

  try {
    console.log(dater);

    const res = await axios.post("http://localhost:8090/calc", dater);

    console.log(`res data:\n${res.data}`);
  } catch (err) {
    console.log("epic fail lmfao:", err);
  }
}
main();
