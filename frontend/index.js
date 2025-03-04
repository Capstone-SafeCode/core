document.addEventListener("DOMContentLoaded", () => {
  const fileInput = document.getElementById("file");
  const submitBtn = document.getElementById("submitBtn");
  const resultText = document.getElementById("result-text");
  let uploadedFileName = null;

  fileInput.addEventListener("change", async function () {
    if (this.files.length > 0) {
        const file = this.files[0];
        const formData = new FormData();
        formData.append("codeFile", file);

        try {
            const response = await fetch("http://localhost:8069/upload", {
                method: "POST",
                body: formData
            });

            if (!response.ok) throw new Error("Upload failed.");

            uploadedFileName = "filled";
        } catch (error) {
            console.error("Error uploading file:", error);
            alert("Error uploading file.");
        }
    }
});

  submitBtn.addEventListener("click", async function () {
      if (!uploadedFileName) {
          alert("Please upload a file first!");
          return;
      }

      try {
        const response = await fetch("http://localhost:8069/start_analyse", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ fileName: uploadedFileName }),
        });

        if (!response.ok) throw new Error("Analysis request failed.");
        const result = await response.json();

        let formattedMessage = result.map(err =>
            `CWE-${err.CWE} detected in ${err.Path} on line [a fix]:
        Kind of vulnerability: ${err.Kind}.
        To solve it: ${err.ToFixIt}`
        ).join("\n\n");

        resultText.textContent = formattedMessage;
    } catch (error) {
        console.error("Error starting analysis:", error);
        alert("Error starting analysis.");
    }
  });
});
