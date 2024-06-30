<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiClient } from "../../http/httpclient";
import { IssueItem } from "../../http/models";

const issues = ref<IssueItem[] | undefined>()

onMounted(async () => {
    const { data } = await apiClient.get<IssueItem[]>("/api/issues")
    issues.value = data
})

</script>

<template>
    <div class="bg-white">
        <div class="py-3 hover:bg-gray-100" v-for="item in issues" v-bind:key="item.id">
            <div class="px-5 text-lg">
                <router-link class="hover:underline underline-offset-2" :to="{ name: 'issue', params: {id: item.id } }">
                    {{ item.title }}
                </router-link>
            </div>
            <div class="px-5 text-sm">
                <span>创建于 {{ item.created_at }}</span>
                <span class="ml-3" v-if="item.created_at != item.updated_at">最后更新于 {{ item.updated_at }}</span>
            </div>
            <div class="px-5">{{ item.description }}</div>

        </div>
    </div>
</template>
