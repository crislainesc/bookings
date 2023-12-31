let attention = Prompt()

;(() => {
  'use strict'

  // Fetch all the forms we want to apply custom Bootstrap validation styles to
  const forms = document.querySelectorAll('.needs-validation')

  // Loop over them and prevent submission
  Array.from(forms).forEach((form) => {
    form.addEventListener(
      'submit',
      (event) => {
        if (!form.checkValidity()) {
          event.preventDefault()
          event.stopPropagation()
        }

        form.classList.add('was-validated')
      },
      false
    )
  })
})()

function notify(msg, type) {
  notie.alert({
    type: type,
    text: msg,
  })
}

function notifyModal({ icon, title, text, confirmationButton, ...rest }) {
  Swal.fire({
    icon,
    title,
    text,
    confirmationButton,
    ...rest,
  })
}

function Prompt() {
  let toast = function (c) {
    const { msg = '', icon = 'success', position = 'top-end' } = c

    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer)
        toast.addEventListener('mouseleave', Swal.resumeTimer)
      },
    })

    Toast.fire({})
  }

  let success = function (c) {
    const { msg = '', title = '', footer = '' } = c

    Swal.fire({
      icon: 'success',
      title: title,
      text: msg,
      footer: footer,
    })
  }

  let error = function (c) {
    const { msg = '', title = '', footer = '' } = c

    Swal.fire({
      icon: 'error',
      title: title,
      text: msg,
      footer: footer,
    })
  }

  async function custom(c) {
    const {
      icon = '',
      msg = '',
      title = '',
      confirmButtonText = 'OK',
      showConfirmButton = true,
    } = c

    const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      backdrop: true,
      focusConfirm: false,
      showConfirmButton: showConfirmButton,
      showCancelButton: true,
      confirmButtonText: confirmButtonText,
      willOpen: () => {
        if (c.willOpen) c.willOpen()
      },
      didOpen: () => {
        if (c.didOpen) c.didOpen()
      },
      preConfirm: () => {
        if (c.preConfirm) c.preConfirm()
      },
      confirm: () => {
        if (c.confirm) c.confirm()
      },
    })

    if (!!result) {
      if (!Swal.DismissReason.cancel) {
        if (!!result.value) {
          if (!!c.callback) {
            c.callback(result)
          }
        } else {
          c.callback(false)
        }
      } else {
        c.callback(true)
      }
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  }
}
