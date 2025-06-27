console.log("works");


async function redirectToGithubLogin() {
    try {
        const res = await fetch("/api/login");
        const resJson = await res.json();
        if (!resJson.ok) {
            throw new Error(resJson.message);
        }
        document.location = resJson.url;
    } catch (error) {
        console.error(error.message);
    }
}
