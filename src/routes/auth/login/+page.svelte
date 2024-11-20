<script lang="ts">
    import { goto } from "$app/navigation";
    import { login } from "$lib/auth";
    import Fa from "svelte-fa";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";
    import { toaster } from "../../../components/toaster/toaster";
    import { faGithub, faGoogle } from "@fortawesome/free-brands-svg-icons";
    import { UserRole } from "$lib/models";

    let email: string;
    let password: string;

    const handleLogin = async () => {
        await login(email, password)
            .then(user => user?.role == UserRole.ADMIN ? goto("/admin") : goto("/user"))
            .catch(() => toaster.push({ type: "error", message: "Failed to login." }, 5000));
    };

</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleLogin} aria-labelledby="login-form">
        <Input bind:value={email} name="Email" id="email" aira-label="Email address"/>
        <Input bind:value={password} name="Password" type="password" id="password" aria-label="Password"/>
        <div class="flex flex-row justify-between">
            <div class="flex flex-row gap-1">
                <Button
                    type="button"
                    on:click={() => window.location.href = "http://localhost:8080/api/auth/github"}
                    aria-label="Sign in with GitHub">
                    <Fa icon={faGithub}/>
                </Button>
                <Button
                    type="button"
                    on:click={() => window.location.href = "http://localhost:8080/api/auth/google"}
                    aria-label="Sign in with Google">
                    <Fa icon={faGoogle}/>
                </Button>
                <!-- <Button type="button" on:click={() => window.location.href = "http://localhost:8080/api/auth/microsoft"}><Fa icon={faMicrosoft}/></Button> -->
            </div>
            <Button type="submit" aria-label="Login">Login</Button>
        </div>
    </form>
</div>
