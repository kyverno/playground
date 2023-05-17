<template>
<v-container fluid>
    <v-row>
        <v-col md="3" cols="6"><span class="d-inline-block font-weight-bold" style="width: 50px">Policy</span>
            {{ details.policy }}</v-col>
        <v-col md="3" cols="6"><span class="d-inline-block font-weight-bold" style="width: 120px">APIVersion/Kind</span>
            {{ details.apiVersion }}/{{ details.kind }}</v-col>
        <v-col cols="3" class="d-none d-md-flex align-center"><span class="d-inline-block font-weight-bold" style="width: 50px">Status</span>
            <StatusChip :status="details.status" />
        </v-col>
        </v-row>
        <v-row class="mt-0">
        <v-col md="3" cols="6"><span class="d-inline-block font-weight-bold" style="width: 50px">Rule</span>{{ details.rule }}</v-col>
        <v-col md="3" cols="6"><span class="d-inline-block font-weight-bold" style="width: 120px">Resource</span>{{ details.resource }}</v-col>
        </v-row>
        <v-row class="d-md-none d-flex">
        <v-col>
            <span class="d-inline-block font-weight-bold" style="width: 50px">Status</span>
            <StatusChip :status="details.status" />
        </v-col>
    </v-row>
</v-container>
</template>

<script setup lang="ts">
import { PropType } from 'vue';
import { RuleStatus } from '@/types';

type Details = {
  apiVersion: string;
  kind: string;
  resource: string;
  policy: string;
  rule: string;
  status: RuleStatus;
};


defineProps({
    details: { type: Object as PropType<Details>, required: true }
})
</script>