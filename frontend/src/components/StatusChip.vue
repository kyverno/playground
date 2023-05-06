<template>
  <v-chip :color="color" dark>
    {{ status }}
  </v-chip>
</template>

<script setup lang="ts">
import { mapStatus } from '../mapper'
import { PropType } from 'vue';
import { RuleStatus } from '@/types';

const props = defineProps({
  status: { type: String as PropType<RuleStatus>, required: true }
})

const color = mapStatus(props.status)
</script>
