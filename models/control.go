package models

import ()

type Control struct {
	Id          int
	Req_id      int
	Cis_id      string
	Category    string
	Requirement string
	Discussion  string
	Check_text  string
	Fix_text    string
	Row_desc    string
	Baseline_id int
}

func (c *Control) SetId(id int) {
	c.Id = id
}

func (c *Control) SetReq_id(req_id int) {
	c.Req_id = req_id
}

func (c *Control) SetBaseline_id(baseline_id int) {
	c.Baseline_id = baseline_id
}

func (c *Control) SetCis_id(cis_id string) {
	c.Cis_id = cis_id
}

func (c *Control) SetCategory(category string) {
	c.Category = category
}

func (c *Control) SetRequirement(requirement string) {
	c.Requirement = requirement
}

func (c *Control) SetDiscussion(discussion string) {
	c.Discussion = discussion
}

func (c *Control) SetCheck_text(check_text string) {
	c.Check_text = check_text
}

func (c *Control) SetFix_text(fix_text string) {
	c.Fix_text = fix_text
}

func (c *Control) SetRow_desc(row_desc string) {
	c.Row_desc = row_desc
}
