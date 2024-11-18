<script lang="ts">
  import { Paginator } from "$lib/api";
  import { authStore } from "$lib/auth";
  import type { Invite } from "$lib/models";
  import { faCancel, faCopy } from "@fortawesome/free-solid-svg-icons";
  import { onMount } from "svelte";
  import { get } from "$lib/api";
  import { getAuthHeader } from "$lib/auth";
  import { IconType } from "$lib/utils";
  import Fa from "svelte-fa";
  import Button from "../../../components/form/Button.svelte";
  import Modal from "../../../components/Modal.svelte";
  import Table from "../../../components/table/Table.svelte";
  import TableItem from "../../../components/table/TableItem.svelte";
  import { toaster } from "../../../components/toaster/toaster";

  let loading = true;
  let data: Invite[] = [];
  let showModal = false;
  let inviteCode = '';

  const paginator = new Paginator<Invite>(
    {
      route: "organization/invites",
      method: "GET",
      headers: getAuthHeader(),
    },
    10,
  );
  paginator.subscribe((items) => (data = items));

  const copyInvite = () => {
    navigator.clipboard.writeText(inviteCode);
    toaster.push({ type: "success", message: `Code "${inviteCode}" copied` }, 5000);
  }

  const createInvite = async () => {
    const response = await get<{ code: string }>({
        route: "organization/invites",
        headers: getAuthHeader(),
        method: "POST",
    });
    if (!response) throw "Failed to create spot.";
    inviteCode = response.code;
  }

  const handleSpotCreation = async () => {
    try {
        await createInvite();
        showModal = true;
    } catch (error) {
        toaster.push({ type: "error", message: "Failed to create invite." }, 5000);
    }
};

  onMount(async () => {
    await paginator.load();
    loading = false;
  });
</script>

<Modal
  visible={showModal}
  title={"Invite Code Created"}
  on:close={() => showModal=false}
  icon = {IconType.SUCCESS}
>
  <h2 id="invite-modal" class="sr-only">Invite Code Created</h2>
  <div id="invite-modal-message" slot="message" class="mt-2 flex" >
    <button class="text-3xl text-gray-900 hover:text-gray-500 font-bold" aria-label="Copy invite code" on:click={() => copyInvite()}>{inviteCode}</button>
    <button class="ml-2 text-gray-500 hover:text-gray-700 focus:outline-none" aria-label="Copy invite code" on:click={() => copyInvite()}>
        <Fa icon={faCopy} class="h-5 w-5"/>
    </button>
  </div>
</Modal>
<Table
  columns={["code", "expiration", "cancel"]}
  {data}
  {loading}
  addRowItem={"invite"}
  addRowFunctionality={() => handleSpotCreation()}
  let:code
  let:expiration
>
  <TableItem>{code}</TableItem>
  <TableItem>{new Date(expiration).toLocaleString()}</TableItem>
  <TableItem
    ><button aria-label={`Cancel invite code ${code}`} on:click={() => console.log("Cancel")}
      ><Fa icon={faCancel} /></button
    ></TableItem
  >
</Table>
<div class="flex flex-row justify-center gap-2">
  <Button on:click={() => paginator.last()} aria-label="Go to the last page">Last</Button>
  <Button on:click={() => paginator.next()} aria-label="Go to the next page">Next</Button>
</div>
