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

  const elem = document.getElementById('reservation-dates')
  new DateRangePicker(elem, {
    format: 'yyyy-mm-dd',
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
    const { msg = '', title = '' } = c

    const { value: formValues } = await Swal.fire({
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      preConfirm: () => {
        return [
          document.getElementById('swal-input1').value,
          document.getElementById('swal-input2').value,
        ]
      },
    })

    if (formValues) {
      Swal.fire(JSON.stringify(formValues))
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  }
}
