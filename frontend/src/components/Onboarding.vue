<template>
  <VOnboardingWrapper ref="wrapper" :steps="steps" @finish="emit('finish')">
    <template #default="{ previous, next, step, isFirst, isLast }">
      <VOnboardingStep>
        <v-card
          max-width="600"
          v-if="step && step.content"
          :title="step.content.title"
          :min-width="350"
          :color="layoutTheme === 'dark' ? 'grey-darken-2' : 'grey-lighten-4'">
          <v-card-text>{{ step.content.description }}</v-card-text>
          <v-card-actions>
            <v-btn @click="close">Close</v-btn>
            <v-btn @click="previous" v-if="!isFirst">Previous</v-btn>
            <v-spacer />
            <v-btn @click="next" v-if="isLast">Finish</v-btn>
            <v-btn @click="next" v-else>Next</v-btn>
          </v-card-actions>
        </v-card>
      </VOnboardingStep>
    </template>
  </VOnboardingWrapper>
</template>

<script setup lang="ts">
import { layoutTheme } from '@/config'
import { VOnboardingWrapper, VOnboardingStep } from 'v-onboarding'
import { watch, ref } from 'vue'

defineProps({
  close: { type: Function },
  steps: { type: Array }
})

const wrapper = ref(null)

watch(wrapper, (w) => emit('wrapper', w))

const emit = defineEmits(['finish', 'wrapper'])
</script>
