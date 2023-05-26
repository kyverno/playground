<template>
  <v-chip :color="color" class="text-white elevation-0" variant="elevated">
    {{ status }}
  </v-chip>
</template>

<script setup lang="ts">
import { PropType } from 'vue'
import { mapStatus } from '@/mapper'
import { RuleStatus } from '@/types'

const props = defineProps({
  status: { type: String as PropType<RuleStatus>, required: true }
})

const color = mapStatus(props.status)
</script>
