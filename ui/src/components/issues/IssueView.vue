<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiClient } from "../../http/httpclient";
import { useRoute } from 'vue-router'
import { IssueItem, ItemSection, SqlQuery as SqlQuery, QueryResult } from "../../http/models";

const route = useRoute()
const issue = ref<IssueItem>()
const sections = ref<ItemSection[]>([])

type NewQueryRequest = {
    sectionId: number
    connectionId: string
    title: string
    query: string
}
const newQueryRequest = ref<NewQueryRequest>()

const issueId = route.params.id

onMounted(async () => {
    const { data: issueData } = await apiClient.get<IssueItem>(`/api/issues/${issueId}`)
    issue.value = issueData

    const { data: sectionsData } = await apiClient.get<ItemSection[]>(`/api/issues/${issueId}/sections`)
    sections.value = sectionsData
})

const onAddQuery = async(sectionid: number) => {
    newQueryRequest.value = {
        sectionId: sectionid,
        connectionId: "maxcompute",
        title: "",
        query: ""
    }
}

const onCreateQuery = async() => {
    const input = newQueryRequest.value!!

    const { data } = await apiClient.post<SqlQuery>(`/api/issues/${issueId}/sections/${newQueryRequest.value?.sectionId}/queries`, {
        connection_id: input.connectionId,
        title: input.title,
        query: input.query
    })

    const currentSections = sections.value
    const section = currentSections.find(x=> x.id == input.sectionId)!!

    section.queries = section.queries || []
    section.queries.push(data)
    return false
}

</script>

<template>
    <div>
        <div class="" v-if="issue">
            <div class="my-3 py-3 bg-white">
                <h1 class="text-2xl	px-5">{{  issue.title  }}</h1>
                <div class="px-5 text-sm">
                    <span>创建于 {{ issue.created_at }}</span>
                    <span class="ml-3" v-if="issue.created_at != issue.updated_at">最后更新于 {{ issue.updated_at }}</span>
                </div>
                <div class="px-5">{{ issue.description }}</div>
            </div>

            <div v-for="section in sections" v-bind:key="section.id" class="bg-white my-3 py-2">
                <div class="m-3">{{ section.header }}</div>
                <div class="m-3">{{ section.body }}</div>

                <div v-for="query in section.queries" v-bind:key="query.id" class="m-3 border p-3">
                    <div>{{ query.query }}</div>
                    <div>{{ query.params }}</div>
                    <div>{{ query.sql }}</div>
                     <table class="border-solid border border-gray-200" v-if="query.result">
                        <thead>
                            <tr>
                                <th class="border px-2" v-for="c in query.result.column_names" v-bind:key="c">{{ c }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr class="border-solid border" v-for="record, idx in query.result.records" v-bind:key="idx">
                                <td class="border px-2" v-for="value in record" v-bind:key="value">{{ value }}</td>
                            </tr>
                        </tbody>
                     </table>
                </div>
                <div class="m-3">{{ section.footer }}</div>

                <div class="my-2 m-3">
                    <button class="bg-orange-500" @click="onAddQuery(section.id)">Add query</button>
                </div>

                <form class="m-3" v-if="newQueryRequest?.sectionId == section.id">
                    <div>
                        <input type="text" v-model="newQueryRequest.connectionId">
                    </div>

                    <div>
                        <input type="text" v-model="newQueryRequest.title">
                    </div>

                    <div>
                        <textarea v-model="newQueryRequest.query" />
                    </div>

                    <div>
                        <button @click="onCreateQuery">Create</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>
