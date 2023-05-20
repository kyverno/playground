import { ContextTemplate, PolicyTemplate, ResourceTemplate } from "@/assets/templates";
import { reactive } from "vue";
import { useState, Inputs } from "@/composables";

export const inputs = reactive({
    policy: PolicyTemplate,
    resource: ResourceTemplate,
    context: ContextTemplate,
    config: '',
})

export const reset = () => {
    inputs.policy = PolicyTemplate
    inputs.resource = ResourceTemplate
    inputs.context = ContextTemplate
}

export const setDefaults = () => init({ policy: PolicyTemplate, resource: ResourceTemplate, context: ContextTemplate, })

export const init = (values: Inputs) => {
    const state = useState()

    if (values.policy) {
        state.policy.value = values.policy;
        inputs.policy = values.policy;
    }

    if (values.resource) {
        state.resource.value = values.resource;
        inputs.resource = values.resource;
    }

    if (values.context) {
        state.context.value = values.context;
        inputs.context = values.context;
    }
    
    state.name.value = values.name || "";
}

export const update = (values: Inputs) => {
    if (values.policy) {
        inputs.policy = values.policy;
    }

    if (values.resource) {
        inputs.resource = values.resource;
    }

    if (values.context) {
        inputs.context = values.context;
    }
}

export const populate = () => {
    const state = useState()

    if (state.policy.value) {
        inputs.policy = state.policy.value;
    }
    if (state.resource.value) {
        inputs.resource = state.resource.value;
    }
    if (state.context.value) {
        inputs.context = state.context.value;
    }
    if (state.config.value) {
        inputs.config = state.config.value;
    }
}