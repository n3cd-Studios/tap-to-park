<script lang="ts">
    import { get } from "$lib/api";
    import { store } from "$lib/auth";
    import Button from "../../components/Button.svelte";
    import Input from "../../components/Input.svelte";

    let email: string;
    let password: string;
    
    const login = async () => {
        const response = await get<{ token: string }>({ route: "auth/login", method: "POST", body: { email, password } });
        if (!response) {
            // failed to login
            return;
        }
        $store = response;
    }

</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <div class="p-10 bg-white rounded-xl w-1/3">
        <Input bind:value={email} name="Email" />
        <Input bind:value={password} name="Password" />
        <div class="flex flex-row justify-end">
            <Button click={login}>Login</Button>
        </div>
    </div>
</div>