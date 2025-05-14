import apiClient from '@/api/client'
import groupsApi from '@/api/endpoints/groups'

export default {
  namespaced: true,

  state: {
    groups: [],
    currentGroup: null,
    members: [],
    isLoading: false,
    error: null,
  },

  getters: {
    allGroups: (state) => state.groups,
    currentGroup: (state) => state.currentGroup,
    groupMembers: (state) => state.members,
    isLoading: (state) => state.isLoading,
    error: (state) => state.error,
  },

  mutations: {
    SET_GROUPS(state, groups) {
      state.groups = groups
    },

    SET_CURRENT_GROUP(state, group) {
      state.currentGroup = group
    },

    SET_GROUP_MEMBERS(state, members) {
      state.members = members
    },

    ADD_GROUP(state, group) {
      state.groups.push(group)
    },

    UPDATE_GROUP(state, updatedGroup) {
      const index = state.groups.findIndex((g) => g.id === updatedGroup.id)
      if (index !== -1) {
        state.groups.splice(index, 1, updatedGroup)
      }

      // Also update current group if it's the same one
      if (state.currentGroup && state.currentGroup.id === updatedGroup.id) {
        state.currentGroup = updatedGroup
      }
    },

    REMOVE_GROUP(state, groupId) {
      state.groups = state.groups.filter((g) => g.id !== groupId)

      // Clear current group if it was the one removed
      if (state.currentGroup && state.currentGroup.id === groupId) {
        state.currentGroup = null
      }
    },

    ADD_MEMBER(state, member) {
      state.members.push(member)
    },

    REMOVE_MEMBER(state, memberId) {
      state.members = state.members.filter((m) => m.id !== memberId)
    },

    SET_LOADING(state, isLoading) {
      state.isLoading = isLoading
    },

    SET_ERROR(state, error) {
      state.error = error
    },
  },

  actions: {
    // Fetch all groups for the current user
    async fetchGroups({ commit }) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        const response = await apiClient.get('/groups')
        commit('SET_GROUPS', response.data)
        return response.data
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to fetch groups')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Fetch a single group by ID
    async fetchGroup({ commit }, groupId) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        const response = await apiClient.get(`/groups/${groupId}`)
        commit('SET_CURRENT_GROUP', response.data)
        return response.data
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to fetch group')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Fetch members of a group
    async fetchGroupMembers({ commit }, groupId) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        const response = await apiClient.get(`/groups/${groupId}/members`)
        commit('SET_GROUP_MEMBERS', response.data)
        return response.data
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to fetch group members')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Create a new group
    async createGroup({ commit }, groupData) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        const response = await apiClient.post('/groups', groupData)
        commit('ADD_GROUP', response.data)
        return response.data
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to create group')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Update an existing group
    async updateGroup({ commit }, { groupId, data }) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        const response = await apiClient.put(`/groups/${groupId}`, data)
        commit('UPDATE_GROUP', response.data)
        return response.data
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to update group')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Delete a group
    async deleteGroup({ commit }, groupId) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        await apiClient.delete(`/groups/${groupId}`)
        commit('REMOVE_GROUP', groupId)
        return true
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to delete group')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Add a member to a group
    async addGroupMember({ commit }, { groupId, userId, role = 'member' }) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        const response = await apiClient.post(`/groups/${groupId}/members`, {
          user_id: userId,
          role,
        })

        commit('ADD_MEMBER', response.data)
        return response.data
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to add member to group')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Remove a member from a group
    async removeGroupMember({ commit }, { groupId, memberId }) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        await apiClient.delete(`/groups/${groupId}/members/${memberId}`)
        commit('REMOVE_MEMBER', memberId)
        return true
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to remove member from group')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // Update a member's role in a group
    async updateMemberRole({ commit, dispatch }, { groupId, memberId, role }) {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)

      try {
        await apiClient.put(`/groups/${groupId}/members/${memberId}`, { role })

        // Refresh the member list to get updated data
        await dispatch('fetchGroupMembers', groupId)

        return true
      } catch (error) {
        commit('SET_ERROR', error.message || 'Failed to update member role')
        throw error
      } finally {
        commit('SET_LOADING', false)
      }
    },
  },
}
