<script setup lang="ts">
import { ref } from "vue";
import { apiClient } from "../../http/httpclient";


const title = ref("");
const description = ref("");

const onCreateIssue = async () => {
    const request = {
        title: title.value,
        description: description.value
    }

    try {
        const { data } = await apiClient.post("/api/issues", request)
        console.log(data)
    } catch(e) {
        console.log(e)
    }
}

</script>

<template>
    <div>
        <h1>New Issue</h1>

        <form>
            <div>
                <input type="text" v-model="title" />
            </div>

            <div>
                <textarea v-model="description" />
            </div>

            <div>
                <button @click="onCreateIssue">Create</button>
            </div>
        </form>
    </div>
</template>
