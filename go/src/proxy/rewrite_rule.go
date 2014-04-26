package proxy

import "strings"
import "strconv"
import "math/rand"
import "net/http"


type RewriteRule interface {
	Active(*http.Request) bool  // check if applies
	Rewrite(*http.Request)      // transform request
}


type RuleTrigger func(*http.Request) bool;


type RewriteRandomRule struct {
	Rules []RewriteRule
}

func (self *RewriteRandomRule) Active(*http.Request) bool {
	return true
}

func (self *RewriteRandomRule) Rewrite(req *http.Request) {
	self.Rules[rand.Intn(len(self.Rules))].Rewrite(req)
}

func NewRewriteRandomRule(rules []RewriteRule) RewriteRule {
	return &RewriteRandomRule{Rules: rules}
}


type RewriteIfRule struct {
	trigger RuleTrigger
	Rule RewriteRule
}

func (self *RewriteIfRule) Active(req *http.Request) bool {
	return self.trigger(req)
}

func (self *RewriteIfRule) Rewrite(req *http.Request) {
	self.Rule.Rewrite(req)
}

func NewRewriteIfRule(trigger RuleTrigger, rule RewriteRule) RewriteRule {
	return &RewriteIfRule{Rule: rule, trigger: trigger}
}


type RewriteHostRule struct {
	Host string
}

func (self *RewriteHostRule) Active(*http.Request) bool {
	return true
}

func (self *RewriteHostRule) Rewrite(req *http.Request) {
	req.Host = self.Host
}

func NewRewriteHostRule(host string) RewriteRule {
	return &RewriteHostRule{Host: host}
}


type RewritePortRule struct {
	Port int
}

func (self *RewritePortRule) Active(*http.Request) bool {
	return true
}

func (self *RewritePortRule) Rewrite(req *http.Request) {
	req.Host = strings.Split(req.Host, ":")[0] + ":" + strconv.Itoa(self.Port)
}

func NewRewritePortRule(port int) RewriteRule {
	return &RewritePortRule{Port: port}
}



