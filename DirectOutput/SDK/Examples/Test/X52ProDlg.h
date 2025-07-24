#pragma once
#include "resource.h"
#include "DirectOutputImpl.h"

///
/// CX52ProDlg
/// 
/// Handles the X52Pro Test Dialog
///
class CX52ProDlg : public CDialogImpl<CX52ProDlg>
{
public:
	CX52ProDlg(CDirectOutput& directoutput, void* hDevice);

	enum { IDD = IDD_X52PRO };

	BEGIN_MSG_MAP(CX52ProDlg)
		MESSAGE_HANDLER(WM_INITDIALOG, OnInitDialog)
		COMMAND_ID_HANDLER(IDOK, OnClose)
		COMMAND_ID_HANDLER(IDCANCEL, OnClose)
		COMMAND_CODE_HANDLER(EN_CHANGE, OnEditChanged)
		COMMAND_CODE_HANDLER(BN_CLICKED, OnCheckChanged)
	END_MSG_MAP()

	LRESULT OnInitDialog(UINT /*uMsg*/, WPARAM /*wParam*/, LPARAM /*lParam*/, BOOL& /*bHandled*/);
	LRESULT OnClose(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnEditChanged(WORD /*wNotifyCode*/, WORD wID, HWND hWndCtl, BOOL& /*bHandled*/);
	LRESULT OnCheckChanged(WORD /*wNotifyCode*/, WORD wID, HWND hWndCtl, BOOL& /*bHandled*/);

	static void __stdcall OnPageChanged(void* hDevice, DWORD dwPage, bool bSetActive, void* pCtxt);
	static void __stdcall OnSoftButtonChanged(void* hDevice, DWORD dwButtons, void* pCtxt);

private:
	CDirectOutput&	m_directoutput;
	void*			m_device;
	int				m_scrollpos;

	void DisplayErrorMessage(LPCTSTR tszMsg, HRESULT hr);
};
