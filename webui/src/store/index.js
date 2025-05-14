import { createStore } from 'vuex'
import auth from './auth'
import conversations from './conversations'
import messages from './messages'
import groups from './groups'

export default createStore({
  strict: process.env.NODE_ENV !== 'production',

  modules: {
    auth,
    conversations,
    messages,
    groups,
  },
})
